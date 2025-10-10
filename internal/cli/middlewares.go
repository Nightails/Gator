package cli

import (
	"context"

	"github.com/Nightails/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		// get current user
		user, err := s.db.GetUserByName(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}
		if err := handler(s, cmd, user); err != nil {
			return err
		}
		return nil
	}
}
