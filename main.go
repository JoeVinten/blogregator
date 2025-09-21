package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JoeVinten/blogregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.ReadConfig()

	if err != nil {
		log.Fatalf("error reading the config: %v", err)
	}

	s := &state{}

	s.cfg = &cfg

	cmdsMap := commands{}

	cmdsMap.handlers = make(map[string]func(*state, command) error)

	cmdsMap.register("login", handlerLogin)

	args := os.Args

	if len(args) < 3 {
		log.Fatalf("not enough arguments passed")
	}

	cmd := command{}

	cmd.name = args[1]
	cmd.arguments = args[2:]

	cmdsMap.run(s, cmd)

	cfg, err = config.ReadConfig()

	if err != nil {
		log.Fatalf("error reading the config: %v", err)
	}

	fmt.Printf("update successful: %+v\n", cfg)
}
