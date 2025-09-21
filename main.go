package main

import (
	"errors"
	"fmt"
	"log"

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

	s.cfg.SetUser(cmd.arguments[0])

	fmt.Printf("username %s, has been set", cmd.arguments[0])

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

	err = cfg.SetUser("joe_vinten")

	if err != nil {
		log.Fatalf("couldn't set current user: %v")
	}

	cfg, err = config.ReadConfig()

	if err != nil {
		log.Fatalf("error reading the config: %v", err)
	}

	fmt.Printf("Update successful: %+v\n", cfg)
}
