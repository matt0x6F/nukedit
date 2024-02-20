package reddit

import (
	"context"
	"os"
	"time"

	"github.com/go-loremipsum/loremipsum"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Client struct {
	client     *reddit.Client
	loremipsum *loremipsum.LoremIpsum
	logger     zerolog.Logger
	debug      bool
}

type NukeRequest struct {
	Scheduled             bool
	CronExpression        string
	Posts                 bool
	Comments              bool
	MaxAge                int
	UseMaxAge             bool
	MaxScore              int
	UseMaxScore           bool
	ReplacementTextLength int
}

// NewClient creates a new client with the given credentials
func NewClient(clientID, clientSecret, username, password string) (*Client, error) {
	credentials := reddit.Credentials{ID: clientID, Secret: clientSecret, Username: username, Password: password}
	client, err := reddit.NewClient(credentials)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:     client,
		loremipsum: loremipsum.New(),
		logger:     log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}, nil
}

// SetDebug sets the debug mode for the client
func (c *Client) SetDebug(debug bool) {
	c.debug = debug

	if debug {
		c.logger = c.logger.Level(zerolog.DebugLevel)
	} else {
		c.logger = c.logger.Level(zerolog.InfoLevel)
	}
}

func (c *Client) Nuke(req NukeRequest) error {
	if req.Posts {
		err := c.EditAndDeleteAllUserPosts(req.ReplacementTextLength)
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to delete posts")
			return err
		}
	}

	if req.Comments {
		err := c.EditAndDeleteAllUserComments(req.ReplacementTextLength)
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to delete comments")
			return err
		}
	}

	return nil
}

// EditAndDeleteAllUserPosts edits and deletes all posts for the authenticated user. It takes a parameter of count to set the number of words for the post.
func (c *Client) EditAndDeleteAllUserPosts(count int) error {
	posts, err := c.RetrieveAllUserPosts()
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to retrieve posts")
		return err
	}

	deleted := 0

	for _, post := range posts {
		err := c.EditPost(post, c.GenerateWords(count))
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to edit post")
			return err
		}

		time.Sleep(1 * time.Second)

		err = c.DeletePost(post)
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to delete post")
			return err
		}

		time.Sleep(1 * time.Second)

		deleted++
	}

	c.logger.Info().Int("PostsDeleted", deleted).Msg("All posts edited and deleted")

	return nil
}

// RetrieveAllUserPosts retrieves all posts for the authenticated user
func (c *Client) RetrieveAllUserPosts() ([]*reddit.Post, error) {
	var (
		posts []*reddit.Post
		after string
	)

	ctx := context.Background()

	for {
		pagePosts, resp, err := c.client.User.Posts(ctx, &reddit.ListUserOverviewOptions{
			ListOptions: reddit.ListOptions{
				Limit: 100,
				After: after,
			},
			Sort: "new",
			Time: "all",
		})
		if err != nil {
			return nil, err
		}

		after = pagePosts[len(pagePosts)-1].FullID
		remaining := resp.Header.Get("X-Ratelimit-Remaining")
		secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

		if c.debug {
			for i, post := range pagePosts {
				c.logger.Debug().Str("PostID", post.FullID).Int("index", i).Str("Created", post.Created.String()).Msg("Post fetched")
			}
		}

		posts = append(posts, pagePosts...)

		c.logger.Debug().Str("RateLimitUsed", resp.Header.Get("X-Ratelimit-Used")).Str("RateLimitRemaining", remaining).Str("RateLimitReset", secUntilReset).Int("PostsCount", len(posts)).Int("PagePosts", len(pagePosts)).Msg("Posts fetched")

		// less than the limit means we're done
		if len(pagePosts) < 100 {
			break
		}

		// if we're out of requests, sleep
		err = evaluateLimit(remaining, secUntilReset)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

// DeletePost deletes a single post by extract its ID
func (c *Client) DeletePost(post *reddit.Post) error {
	ctx := context.Background()

	resp, err := c.client.Post.Delete(ctx, post.FullID)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	c.logger.Debug().Str("RateLimitUsed", resp.Header.Get("X-Ratelimit-Used")).Str("RateLimitRemaining", remaining).Str("RateLimitReset", secUntilReset).Msg("Post deleted")

	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// EditPost edits a single post, setting the content to the given words
func (c *Client) EditPost(post *reddit.Post, words string) error {
	ctx := context.Background()

	_, resp, err := c.client.Post.Edit(ctx, post.FullID, words)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	c.logger.Debug().Str("RateLimitUsed", resp.Header.Get("X-Ratelimit-Used")).Str("RateLimitRemaining", remaining).Str("RateLimitReset", secUntilReset).Msg("Post edited")

	// if we're out of requests, sleep
	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// EditAndDeleteAllUserComments edits and deletes all comments for the authenticated user. It takes a parameter of count to set the number of words for the comment.
func (c *Client) EditAndDeleteAllUserComments(count int) error {
	comments, err := c.RetrieveAllUserComments()
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to retrieve comments")
		return err
	}

	deleted := 0

	for _, comment := range comments {
		err := c.EditComment(comment, c.GenerateWords(count))
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to edit comment")
			return err
		}

		time.Sleep(1 * time.Second)

		err = c.DeleteComment(comment)
		if err != nil {
			c.logger.Error().Err(err).Msg("Failed to delete comment")
			return err
		}

		time.Sleep(1 * time.Second)

		deleted++
	}

	c.logger.Info().Int("CommentsDeleted", deleted).Msg("All comments edited and deleted")

	return nil
}

// GenerateWords generates words using a LoremIpsum generator
func (c *Client) GenerateWords(quantity int) string {
	return c.loremipsum.Words(quantity)
}

// DeleteComment deletes a single comment by extract its ID
func (c *Client) DeleteComment(comment *reddit.Comment) error {
	ctx := context.Background()

	resp, err := c.client.Comment.Delete(ctx, comment.FullID)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	c.logger.Debug().Str("RateLimitUsed", resp.Header.Get("X-Ratelimit-Used")).Str("RateLimitRemaining", remaining).Str("RateLimitReset", secUntilReset).Msg("Comment deleted")

	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// EditComment edits a single comment, setting the content to the given words
func (c *Client) EditComment(comment *reddit.Comment, words string) error {
	ctx := context.Background()

	_, resp, err := c.client.Comment.Edit(ctx, comment.FullID, words)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	c.logger.Debug().Str("RateLimitUsed", resp.Header.Get("X-Ratelimit-Used")).Str("RateLimitRemaining", remaining).Str("RateLimitReset", secUntilReset).Msg("Comment edited")

	// if we're out of requests, sleep
	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// RetrieveAllUserComments retrieves all comments for the authenticated user
func (c *Client) RetrieveAllUserComments() ([]*reddit.Comment, error) {
	var (
		comments []*reddit.Comment
		after    string
	)

	ctx := context.Background()

	for {
		pageComments, resp, err := c.client.User.Comments(ctx, &reddit.ListUserOverviewOptions{
			ListOptions: reddit.ListOptions{
				Limit: 100,
				After: after,
			},
			Sort: "new",
			Time: "all",
		})
		if err != nil {
			return nil, err
		}

		after = pageComments[len(pageComments)-1].FullID
		remaining := resp.Header.Get("X-Ratelimit-Remaining")
		secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

		if c.debug {
			for i, comment := range pageComments {
				c.logger.Debug().Str("CommentID", comment.FullID).Int("index", i).Str("Created", comment.Created.String()).Msg("Comment fetched")
			}
		}

		comments = append(comments, pageComments...)

		c.logger.Debug().Str("RateLimitUsed", resp.Header.Get("X-Ratelimit-Used")).Str("RateLimitRemaining", remaining).Str("RateLimitReset", secUntilReset).Int("CommentsCount", len(comments)).Int("PageComments", len(pageComments)).Msg("Comments fetched")

		// less than the limit means we're done
		if len(pageComments) < 100 {
			break
		}

		// if we're out of requests, sleep
		err = evaluateLimit(remaining, secUntilReset)
		if err != nil {
			return nil, err
		}
	}

	return comments, nil
}
