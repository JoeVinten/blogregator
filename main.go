package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/JoeVinten/blogregator/internal/config"
	"github.com/JoeVinten/blogregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.ReadConfig()

	if err != nil {
		log.Fatalf("error reading the config: %v", err)
	}

	s := &state{}

	s.cfg = &cfg

	db, err := sql.Open("postgres", cfg.DBURL)

	if err != nil {
		log.Fatalf("error, %v opening connection to %s", err, cfg.DBURL)
	}

	dbQueries := database.New(db)

	s.db = dbQueries

	cmdsMap := commands{}

	cmdsMap.handlers = make(map[string]func(*state, command) error)

	cmdsMap.register("login", handlerLogin)

	cmdsMap.register("register", handlerRegister)

	cmdsMap.register("reset", handlerReset)

	args := os.Args

	if len(args) < 2 {
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
