package main

import (
	"fmt"
	"os"

	"github.com/Nightails/gator/internal/command"
	"github.com/Nightails/gator/internal/config"
)

func main() {
	cfg := config.Read()
	_ = command.State{Config: &cfg}
	commands := command.Commands{CmdMap: make(map[string]func(*command.State, command.Command) error)}
	commands.Register("login", command.HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: missing command")
		os.Exit(1)
	}
	cmd := command.Command{Name: args[1], Args: args[2:]}
	if err := commands.Run(&command.State{Config: &cfg}, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
