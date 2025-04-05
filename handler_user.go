package main

import (
	"fmt"
	"errors"
	"context"
	"github.com/FallenL3vi/blogaggregator/internal/database"
	"time"
	"github.com/google/uuid"
	"database/sql"
	"github.com/lib/pq"
	"strconv"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New(fmt.Sprintf("usage: %s <name>", cmd.Name))
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])

	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(cmd.Args[0])

	if err != nil {
		return errors.New(fmt.Sprintf("couldn't set current user: %w", err))
	}

	fmt.Println("User switched successfully!")
	
	return nil

}

func printUser(user database.User) {
	fmt.Printf(" * ID:	%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return errors.New(fmt.Sprintf("usage: %s", cmd.Name))
	}

	err := s.db.DeleteUsers(context.Background())

	if err != nil {
		return fmt.Errorf("couldn't delete users\n")
	}

	fmt.Println("Table users restarted!\n")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return errors.New(fmt.Sprintf("usage: %s", cmd.Name))
	}

	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("couldn't print users: %w", err)
	}

	for _, v := range users {
		if v.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", v.Name)
		} else {
			fmt.Printf("* %s\n", v.Name)
		}
	}

	return nil
}


func handlerAggreg(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return errors.New(fmt.Sprintf("usage: %s [1s/1m/1h]\n", cmd.Name))
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])

	if err != nil {
		return fmt.Errorf("couldn't pare time\n")
	}

	ticker := time.NewTicker(time_between_reqs)

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	for ; ; <-ticker.C {
		fmt.Printf(" * New batch *\n")
		err := scrapFeeds(s)
		if err != nil {
			fmt.Printf("Couldn't scrap the feed: %v\n", err)
			break
		}
	}
	return nil
}

func scrapFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return fmt.Errorf("No feeds to fetch\n")
	}

	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		return fmt.Errorf("Coultn't mark feed as fetched\n")
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch the feed: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}

		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time: t,
				Valid: true,
			}
		}
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: sql.NullString{
				String: item.Title,
				Valid: true,
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			PublishedAt: publishedAt,
			FeedID: feed.ID,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			fmt.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	fmt.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return errors.New(fmt.Sprintf("usage: %s <LIMIT optional>", cmd.Name))
	}

	var limit int = 2

	if len(cmd.Args) != 0 {
		new_limit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Couldn't convert string to int\n")
		}

		limit = new_limit
	}

	posts, err := s.db.GetPosts(context.Background(), database.GetPostsParams{
		Limit: int32(limit),
		UserID: user.ID,
	})

	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf(" * Created at: %s\n", post.CreatedAt)
		fmt.Printf(" * Title: %s\n", post.Title)
		fmt.Printf(" * Description: %s\n", post.Description)
		fmt.Printf("======================================\n")
	}
	return nil
}