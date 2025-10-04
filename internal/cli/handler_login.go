package cli

import (
	"context"
	"errors"
	"fmt"
)

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
