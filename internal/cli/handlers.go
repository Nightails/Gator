package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nightails/gator/internal/database"
	"github.com/Nightails/gator/internal/rss"
	"github.com/google/uuid"
)

// handlerAddFeed adds a new feed for the current user, stores it in the database, and sets the user to follow the feed.
// It validates the command arguments, retrieves the current user from the database, creates a feed, and follows it.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("missing feed name and url")
	}

	ctx := context.Background()

	// create the feed
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return err
	}

	fmt.Printf("added feed: %s\n", feed.Name)
	fmt.Printf("feed url: %s\n", feed.Url)
	fmt.Printf("created at: %v\n", feed.CreatedAt)

	// follow the feed
	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	ffRecord, err := s.db.CreateFeedFollow(ctx, ffParams)
	if err != nil {
		return errors.New("unable to follow this feed")
	}
	fmt.Printf("%s follow %s\n", ffRecord.UserName, ffRecord.FeedName)

	return nil
}

// handlerAgg fetches the RSS feed from the given URL and prints it to the console.
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	// temporary url, to be replaced with config
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	rssf, err := rss.FetchFeed(ctx, url)
	if err != nil {
		return err
	}
	rssf.UnescapeString()
	printRSSFeed(rssf)

	return nil
}

// printRSSFeed prints the given RSS feed to the console.
func printRSSFeed(rssf *rss.RSSFeed) {
	fmt.Printf("Title: %s\n", rssf.Channel.Title)
	fmt.Printf("Description: %s\n", rssf.Channel.Description)
	fmt.Println("Items:")
	for _, item := range rssf.Channel.Item {
		fmt.Printf("- Title: %s\n", item.Title)
		fmt.Printf("- Link: %s\n", item.Link)
		fmt.Printf("- Description: %s\n", item.Description)
		fmt.Println()
	}
}

// handlerFeeds lists all the feeds in the database.
func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	fmt.Println("listing feeds:")
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.Name)
		fmt.Printf("- url: %s\n", feed.Url)
		// get the user that created the feed
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("- created by: %s\n", user.Name)
		fmt.Println()
	}

	return nil
}

// handlerFollow adds a new feed for the current user, stores it in the database, and sets the user to follow the feed.
// It validates the command arguments, retrieves the current user from the database, creates a feed, and follows it.
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("missing url")
	}

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return errors.New("this feed does not exist")
	}
	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	ffRecord, err := s.db.CreateFeedFollow(ctx, ffParams)
	if err != nil {
		return errors.New("unable to follow this feed")
	}

	fmt.Printf("%s follow follow follow %s\n", ffRecord.UserName, ffRecord.FeedName)
	return nil
}

// handlerFollowing lists all the feeds that the current user is following.
func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	ctx := context.Background()

	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return errors.New("unable to retrieve following feeds")
	}

	fmt.Printf("%s following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.FeedName)
	}

	return nil
}

// handlerLogin logs in the user with the given username.
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}

	ctx := context.Background()
	// Check if the user exists in the database
	user, err := s.db.GetUserByName(ctx, cmd.args[0])
	if err != nil {
		return errors.New("user does not exist")
	}

	s.cfg.SetUser(user.Name)
	fmt.Printf("logged in as %s\n", s.cfg.UserName)
	return nil
}

// handlerRegister creates a new user in the database.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}

	ctx := context.Background()

	// Check if the user already exists
	if _, err := s.db.GetUserByName(ctx, cmd.args[0]); err == nil {
		return fmt.Errorf("user %s already exists\n", cmd.args[0])
	}

	// Create a user in the database
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Name:      cmd.args[0],
	}
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	s.cfg.SetUser(user.Name)
	fmt.Printf("registered and logged in as user: %s\n", user.Name)
	return nil
}

// handlerReset removes all the users from the database.
func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	ctx := context.Background()
	if err := s.db.RemoveUsers(ctx); err != nil {
		return errors.New("failed to reset database")
	}

	return nil
}

// handlerUsers lists all the users in the database.
func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return errors.New("failed to get users")
	}
	for _, user := range users {
		if user.Name == s.cfg.UserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("missing feed url")
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return errors.New("failed to get feed")
	}
	if err := s.db.RemoveFeedFollow(ctx, database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return errors.New("failed to remove following feed")
	}
	return nil
}
