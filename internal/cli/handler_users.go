package cli

import (
	"context"
	"errors"
	"fmt"
)

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
