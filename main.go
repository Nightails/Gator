package main

import (
	"database/sql"
	"fmt"

	"github.com/Nightails/gator/internal/cli"
	"github.com/Nightails/gator/internal/config"
	"github.com/Nightails/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		fmt.Printf("error: failed to open database: %v\n", err)
	}
	dbQueries := database.New(db)
	cli.RunCli(&cfg, dbQueries)
}
