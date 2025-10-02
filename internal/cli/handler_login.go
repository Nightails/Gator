package cli

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("Logged in as %s\n", s.cfg.UserName)
	return nil
}
