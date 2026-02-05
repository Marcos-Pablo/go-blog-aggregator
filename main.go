package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Marcos-Pablo/go-blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("Error while reading config file: %v\n", err)
		os.Exit(1)
	}

	programState := &state{
		cfg: &config,
	}

	commands := commands{
		registry: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)
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

	err = commands.run(programState, command)
	if err != nil {
		log.Fatal(err)
	}
}
