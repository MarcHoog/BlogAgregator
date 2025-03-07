package cli

func NewExplorerClient() *Client {

	c := Client{}
	c.Commands = make(map[string]func(*State, Command) error)
	c.Register("login", handlerLogin)
	c.Register("register", handlerRegister)
	c.Register("reset", handlerReset)
	c.Register("users", handlerUsers)
	c.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.Register("feeds", handlerListFeeds)
	c.Register("follow", middlewareLoggedIn(handlerFollowFeed))
	c.Register("following", middlewareLoggedIn(handlerListFollows))
	c.Register("unfollow", middlewareLoggedIn(handlerUnFollowFeed))
	return &c
}

func NewAggregatorClient() *Client {

	c := Client{}
	c.Commands = make(map[string]func(*State, Command) error)
	c.Register("agg", handleAggregator)

	return &c
}
