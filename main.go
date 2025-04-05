
package main

import (
	"log"
	"database/sql"
	"github.com/FallenL3vi/blogaggregator/internal/config"
	"github.com/FallenL3vi/blogaggregator/internal/database"
	_ "github.com/lib/pq"
	"os"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
		
	if err != nil {
		log.Fatalf("error reading config %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)

	if err != nil {
		log.Fatalf("error could not connect to db: %v", err)
	}
	
	defer db.Close()

	dbQueries := database.New(db)

	var currentState *state = &state{db : dbQueries, cfg : &cfg}

	var cmd commands = commands{registeredCommands : make(map[string]func(*state, command) error),}


	cmd.register("login", handlerLogin)
	cmd.register("register", handlerRegister)
	cmd.register("reset", handlerReset)
	cmd.register("users", handlerGetUsers)
	cmd.register("agg", handlerAggreg)
	cmd.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmd.register("feeds", handlerGetFeeds)
	cmd.register("follow", middlewareLoggedIn(handlerFollow))
	cmd.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmd.register("following", handlerFollowing)
	cmd.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]\n")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmd.run(currentState, command{Name : cmdName, Args : cmdArgs,},)

	if err != nil {
		log.Fatal(err)
	}

}