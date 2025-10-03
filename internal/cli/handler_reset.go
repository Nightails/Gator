package cli

import (
	"context"
	"errors"
)

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
