package main

import (
	"fmt"
	"errors"
	"context"
	"github.com/FallenL3vi/blogaggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if s.cfg.CurrentUserName == "" {
			return errors.New("No user logged in\n")
		}

		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)

		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("Wrong number of arrguments: requires 2\n")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	})

	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	printFeed(&feed)

	return nil
}


func printFeed(feed *database.Feed) {
	fmt.Printf(" * ID:	%v\n", feed.ID)
	fmt.Printf(" * CreatedAt:	%v\n", feed.CreatedAt)
	fmt.Printf(" * UpdatedAt:	%v\n", feed.UpdatedAt)
	fmt.Printf(" * Name:	%v\n", feed.Name)
	fmt.Printf(" * URL:	%v\n", feed.Url)
	fmt.Printf(" * UserID:	%v\n", feed.UserID)
}


func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return errors.New(fmt.Sprintf("usage: %s <name>", cmd.Name))
	}

	feeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return nil
	}

	for _, v := range feeds {
		fmt.Printf(" * Name: %v\n", v.Name)
		fmt.Printf(" * URL: %v\n", v.Url)
		fmt.Printf(" * User: %v\n", v.Name_2)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return errors.New(fmt.Sprintf("usage: %s <name>", cmd.Name))
	}

	var url string = cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)

	if err != nil {
		return err
	}

	follows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	for _, v := range follows {
		fmt.Printf(" * Feed name: %v\n", v.FeedName)
		fmt.Printf(" * User name: %v\n", v.UserName)
	}

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return errors.New(fmt.Sprintf("usage: %s", cmd.Name))
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)

	if err != nil {
		return err
	}

	fmt.Printf(" * Following feeds: \n")
	for _, v := range follows {
		fmt.Printf(" * Feed name: %v\n", v.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return errors.New(fmt.Sprintf("usage: %s <url>\n", cmd.Name))
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])

	if err != nil {
		return err
	}

	err = s.db.DeleteFollow(context.Background(), database.DeleteFollowParams{feed.ID, user.ID})

	if err != nil {
		return err
	}

	return nil
}