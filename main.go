package main

import (
	"database/sql"
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

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error, %v opening connection to %s", err, cfg.DBURL)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmdsMap := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmdsMap.register("login", handlerLogin)
	cmdsMap.register("register", handlerRegister)
	cmdsMap.register("reset", handlerReset)
	cmdsMap.register("users", handlerUsers)
	cmdsMap.register("agg", handlerFetchFeed)
	cmdsMap.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmdsMap.register("feeds", handlerListFeeds)
	cmdsMap.register("follow", middlewareLoggedIn(handlerFollow))
	cmdsMap.register("following", middlewareLoggedIn(handlerFollowing))
	cmdsMap.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments passed")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	err = cmdsMap.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatalf("run command failed to run: %s", err)
	}

}
