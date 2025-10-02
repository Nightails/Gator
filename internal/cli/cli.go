package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/Nightails/gator/internal/config"
	"github.com/Nightails/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdMap map[string]func(*state, command) error
}

func RunCli(cfg *config.Config) {
	s, cmds := setupCli(cfg)
	registerCommands(&cmds)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: missing command")
		os.Exit(1)
	}
	cmd := command{name: args[1], args: args[2:]}
	if err := cmds.run(&s, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func setupCli(cfg *config.Config) (state, commands) {
	s := state{cfg: cfg}
	cmds := commands{make(map[string]func(*state, command) error)}
	return s, cmds
}

func registerCommands(cmds *commands) {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmdMap[cmd.name]
	if !ok {
		return errors.New("unknown command")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.cmdMap[name] = handler
}
