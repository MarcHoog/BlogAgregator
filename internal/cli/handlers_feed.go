package cli

import (
	"bootdevBlogAggerator/internal/database"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func handlerAddFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) < 2 {
		return fmt.Errorf("Expected atleast two arguments <feedname> <feedurl>\n")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	paramsFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), paramsFeed)
	if err != nil {
		return fmt.Errorf("Something went wrong creating feed: %v\n", err)
	}

	printFeed(feed)

	paramsFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), paramsFollow)
	if err != nil {
		return fmt.Errorf("Something went wrong creating feed follow: %v\n", err)
	}

	fmt.Printf(" * %v Follows %v \n", feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func handlerListFeeds(s *State, cmd Command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong getting feeds: %v\n", err)
	}

	userCache := make(map[uuid.UUID]database.User)

	for _, feed := range feeds {

		user, ok := userCache[feed.UserID]

		if !ok {
			user, err = s.db.GetUserById(context.Background(), feed.UserID)

			if err != nil {
				return fmt.Errorf("Something went wrong getting user by id: %v\n", err)
			}

			userCache[user.ID] = user

		}

		printFeed(feed)

	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" * URL:     %v\n", feed.Url)
	fmt.Printf(" * Created: %v\n", feed.CreatedAt)
	fmt.Printf(" * Updated: %v\n", feed.UpdatedAt)
	fmt.Printf(" * Last Fetched: %v\n", feed.LastFetchedAt)
	fmt.Println()
}
