package main

import (
	"github.com/Nightails/gator/internal/cli"
	"github.com/Nightails/gator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()
	cli.RunCli(&cfg)
}
