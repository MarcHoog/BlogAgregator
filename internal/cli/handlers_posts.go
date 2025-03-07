package cli

import (
	"bootdevBlogAggerator/internal/database"
	"bootdevBlogAggerator/internal/rss"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func parsePostPubDate(pubDate string) (time.Time, error) {

	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
	}

	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("failed to parse pub date")
}

func createPostsFromFeed(s *State, feed database.Feed, rssFeed *rss.Feed) error {

	for _, item := range rssFeed.Channel.Item {
		exists, err := s.db.CheckPostExists(context.Background(), item.Link)
		fmt.Printf("post with the url %s exists: %v\n", item.Link, exists)
		if err != nil {
			return err
		}

		if exists {
			continue
		}

		pubDate, err := parsePostPubDate(item.PubDate)
		if err != nil {
			fmt.Printf("Error parsing post pub date '%s' \n", item.PubDate)
			continue
		}

		postsParams := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: false},
			FeedID:      feed.ID,
			PublishedAt: pubDate,
		}

		post, err := s.db.CreatePost(context.Background(), postsParams)
		if err != nil {
			return err
		}

		PrintPost(post)
	}
	return nil
}

func scrapeFeed(s *State, feed database.Feed) (*rss.Feed, error) {
	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong fetching feed: %v\n", err)
	}

	rss.CleanFeed(rssFeed)

	lastFetched := sql.NullTime{Time: time.Now().UTC(), Valid: true}
	param := database.MarkFeedFetchedParams{ID: feed.ID, LastFetchedAt: lastFetched}
	err = s.db.MarkFeedFetched(context.Background(), param)
	if err != nil {
		return nil, err
	}

	return rssFeed, nil
}

func handleAggregator(s *State, cmd Command) error {

	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	_, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("bad time between `%s`: %v", cmd.Args[0], err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {

		feed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return err
		}

		rssFeed, err := scrapeFeed(s, feed)
		if err != nil {
			return err
		}

		err = createPostsFromFeed(s, feed, rssFeed)
		if err != nil {
			return err
		}

	}
}

func handleBrowsePosts(s *State, cmd Command, user database.User) error {

	params := database.GetPostsForUserParams{UserID: user.ID, Limit: 5}

	posts, err := s.db.GetPostsForUser(context.Background(), user.ID)

	for _, post := range posts {
		PrintPost(post)
	}
	return nil
}

func PrintPost(post database.Post) {

	fmt.Printf("Post ID: %s\n", post.ID)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Println()

}
