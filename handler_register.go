package main

import (
	"fmt"
	"errors"
	"context"
	"time"
	"github.com/FallenL3vi/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New(fmt.Sprintf("usage: %s <name>", cmd.Name))
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:	uuid.New(),
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		Name:	name,
	})

	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)

	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully!")
	
	printUser(user)
	return nil

}