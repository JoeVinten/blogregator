package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/JoeVinten/blogregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return errors.New("no username given")
	}

	err := s.cfg.SetUser(cmd.arguments[0])

	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("username %s, has been set\n", cmd.arguments[0])

	return nil

}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]

	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
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
