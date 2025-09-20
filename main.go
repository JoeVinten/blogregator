package main

import (
	"fmt"
	"log"

	"github.com/JoeVinten/blogregator/internal/config"
)

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
