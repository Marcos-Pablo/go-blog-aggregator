package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/config"
	"github.com/Marcos-Pablo/go-blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg     *config.Config
	queries *database.Queries
	db      *sql.DB
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", config.DbURL)
	if err != nil {
		log.Fatalf("Error connect to db: %v", err)
	}

	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		cfg:     &config,
		queries: dbQueries,
		db:      db,
	}

	cmds := getCommands()
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := args[1]
	args = args[2:]

	command := command{
		name: cmdName,
		args: args,
	}

	err = cmds.run(programState, command)

	if err != nil {
		log.Fatal(err)
	}
}
