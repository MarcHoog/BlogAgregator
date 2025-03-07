package cli

import (
	"bootdevBlogAggerator/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func handlerFollowFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("Expected Feed url <feedurl>\n")
	}

	feedUrl := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Feed url doesn't exists add feed by using \n")
		} else {
			return fmt.Errorf("Something went wrong fetching feed by using \n")
		}
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Something went wrong creating feed follow: %v\n", err)
	}

	fmt.Printf(" * %v Follows %v \n", feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func handlerListFollows(s *State, cmd Command, user database.User) error {

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Something went wrong getting feed follows: %v\n", err)
	}

	for _, follow := range follows {

		fmt.Printf(" * FeedName: %v\n", follow.FeedName)
		fmt.Printf(" 	* FollowedAt: %v\n", follow.CreatedAt.Format(time.Stamp))

	}

	return nil
}

func handlerUnFollowFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("Expected Feed url <feedurl>\n")
	}

	url := cmd.Args[0]

	param := database.GetFeedFollowByUserAndUrlParams{
		Url:    url,
		UserID: user.ID,
	}

	_, err := s.db.GetFeedFollowByUserAndUrl(context.Background(), param)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("You are not following feed Or Feed doesn't exist %v\n", url)
		return nil
	} else if err != nil {
		fmt.Printf("Something went wrong fetching feed %v\n", url)
	}

	deleteParam := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteParam)
	if err != nil {
		return fmt.Errorf("Something went wrong deleting feed follow: %v\n", err)
	}

	fmt.Printf(" * Unfollowed Feed %v\n", url)

	return nil
}
