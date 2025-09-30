package command

import (
	"errors"
	"fmt"

	"github.com/Nightails/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CmdMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.CmdMap[cmd.Name]
	if !ok {
		return errors.New("unknown Command")
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	c.CmdMap[name] = handler
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("missing username")
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Printf("Logged in as %s\n", s.Config.UserName)
	return nil
}
