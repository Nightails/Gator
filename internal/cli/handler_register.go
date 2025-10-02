package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nightails/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}

	ctx := context.Background()

	// Check if the user already exists
	if _, err := s.db.GetUsers(ctx, cmd.args[0]); err == nil {
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
