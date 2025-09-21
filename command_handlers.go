package main

import (
	"errors"
	"fmt"
)

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
