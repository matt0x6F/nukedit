package reddit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-loremipsum/loremipsum"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Log struct {
	Time    time.Time `json:"timestamp"`
	Message string    `json:"message"`
	Done    bool      `json:"done"`
	Error   error     `json:"error"`
}

type Client struct {
	ctx        context.Context
	client     *reddit.Client
	loremipsum *loremipsum.LoremIpsum
	debug      bool
	dryRun     bool
}

type NukeRequest struct {
	Scheduled             bool   `json:"scheduled"`
	CronExpression        string `json:"cronExpression"`
	Posts                 bool   `json:"posts"`
	Comments              bool   `json:"comments"`
	MaxAge                int    `json:"maxAge"`
	UseMaxAge             bool   `json:"useMaxAge"`
	MinScore              int    `json:"minScore"`
	UseMinScore           bool   `json:"useMinScore"`
	ReplacementTextLength int    `json:"replacementTextLength"`
}

type NukeResult struct {
	CommentsDeleted int `json:"commentsDeleted"`
	PostsDeleted    int `json:"postsDeleted"`
}

// NewClient creates a new client with the given credentials
func NewClient(ctx context.Context, clientID, clientSecret, username, password string, dryRun bool) (*Client, error) {
	credentials := reddit.Credentials{ID: clientID, Secret: clientSecret, Username: username, Password: password}
	client, err := reddit.NewClient(credentials)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:        ctx,
		client:     client,
		loremipsum: loremipsum.New(),
		dryRun:     dryRun,
		debug:      dryRun,
	}, nil
}

// SetDebug sets the debug mode for the client
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) Nuke(req NukeRequest, logChan chan<- Log) (NukeResult, error) {
	result := NukeResult{}

	if req.Posts {
		posts, err := c.EditAndDeleteAllUserPosts(req, logChan)
		result.PostsDeleted = posts
		if err != nil {
			runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to delete posts: %s", err.Error())
			logChan <- Log{Time: time.Now(), Message: "Failed to delete posts", Done: true, Error: err}
			return result, err
		}
	}

	if req.Comments {
		comments, err := c.EditAndDeleteAllUserComments(req, logChan)
		result.CommentsDeleted = comments
		if err != nil {
			runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to delete comments: %s", err.Error())
			logChan <- Log{Time: time.Now(), Message: "Failed to delete comments", Done: true, Error: err}
			return result, err
		}
	}

	logChan <- Log{Time: time.Now(), Message: "Nuke complete", Done: true, Error: nil}

	return result, nil
}

// EditAndDeleteAllUserPosts edits and deletes all posts for the authenticated user. It takes a parameter of count to set the number of words for the post.
func (c *Client) EditAndDeleteAllUserPosts(req NukeRequest, logChan chan<- Log) (int, error) {
	posts, err := c.RetrieveAllUserPosts(logChan)
	if err != nil {
		runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to retrieve posts: %s", err.Error())
		return 0, err
	}

	deleted := 0

	for _, post := range posts {
		// Evaluate request and skip ones that match the criteria
		if req.UseMaxAge && time.Since(post.Created.Time) < time.Duration(req.MaxAge*24)*time.Hour {
			continue
		}

		if req.UseMinScore && post.Score > req.MinScore {
			continue
		}

		err := c.EditPost(post, c.GenerateWords(req.ReplacementTextLength), logChan)
		if err != nil {
			runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to edit post: %s", err.Error())
			return deleted, err
		}

		time.Sleep(1 * time.Second)

		err = c.DeletePost(post, logChan)
		if err != nil {
			runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to delete post: %s", err.Error())
			return deleted, err
		}

		time.Sleep(1 * time.Second)

		deleted++
	}

	runtime.LogInfof(c.ctx, "[Reddit Client] All posts edited and deleted: %d", deleted)

	return deleted, nil
}

// RetrieveAllUserPosts retrieves all posts for the authenticated user
func (c *Client) RetrieveAllUserPosts(logChan chan<- Log) ([]*reddit.Post, error) {
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
			for _, post := range pagePosts {
				runtime.LogDebugf(c.ctx, "[Reddit Client] Post fetched: %s", post.FullID)
				logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Post fetched: %s", post.FullID), Done: false, Error: nil}
			}
		}

		posts = append(posts, pagePosts...)

		runtime.LogDebugf(c.ctx, "[Reddit Client] Posts fetched: %d", len(posts))
		logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Fetched %d posts; total: %d", len(pagePosts), len(posts)), Done: false, Error: nil}

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
func (c *Client) DeletePost(post *reddit.Post, logChan chan<- Log) error {
	ctx := context.Background()

	if c.dryRun {
		runtime.LogInfof(c.ctx, "[Reddit Client] Dry run: Post would be deleted: %s", post.FullID)
		logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Dry run: Post would be deleted: %s", post.FullID), Done: false, Error: nil}
		return nil
	}

	resp, err := c.client.Post.Delete(ctx, post.FullID)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	runtime.LogDebugf(c.ctx, "[Reddit Client] Post deleted: %s", post.FullID)

	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// EditPost edits a single post, setting the content to the given words
func (c *Client) EditPost(post *reddit.Post, words string, logChan chan<- Log) error {
	ctx := context.Background()

	if c.dryRun {
		runtime.LogInfof(c.ctx, "[Reddit Client] Dry run: Post would be edited: %s", post.FullID)
		logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Dry run: Post would be edited: %s", post.FullID), Done: false, Error: nil}
		return nil
	}

	_, resp, err := c.client.Post.Edit(ctx, post.FullID, words)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	runtime.LogDebugf(c.ctx, "[Reddit Client] Post edited: %s", post.FullID)

	// if we're out of requests, sleep
	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// EditAndDeleteAllUserComments edits and deletes all comments for the authenticated user. It takes a parameter of count to set the number of words for the comment.
func (c *Client) EditAndDeleteAllUserComments(req NukeRequest, logChan chan<- Log) (int, error) {
	comments, err := c.RetrieveAllUserComments(logChan)
	if err != nil {
		runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to retrieve comments: %s", err.Error())
		return 0, err
	}

	deleted := 0

	for _, comment := range comments {
		// Evaluate request and skip ones that match the criteria
		if req.UseMaxAge && time.Since(comment.Created.Time) < time.Duration(req.MaxAge*24)*time.Hour {
			continue
		}

		if req.UseMinScore && comment.Score > req.MinScore {
			continue
		}

		err := c.EditComment(comment, c.GenerateWords(req.ReplacementTextLength), logChan)
		if err != nil {
			runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to edit comment: %s", err.Error())
			return deleted, err
		}

		time.Sleep(1 * time.Second)

		err = c.DeleteComment(comment, logChan)
		if err != nil {
			runtime.LogErrorf(c.ctx, "[Reddit Client] Failed to delete comment: %s", err.Error())
			return deleted, err
		}

		time.Sleep(1 * time.Second)

		deleted++
	}

	runtime.LogInfof(c.ctx, "[Reddit Client] All comments edited and deleted: %d", deleted)
	logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("All comments edited and deleted: %d", deleted), Done: false, Error: nil}

	return deleted, nil
}

// GenerateWords generates words using a LoremIpsum generator
func (c *Client) GenerateWords(quantity int) string {
	return c.loremipsum.Words(quantity)
}

// DeleteComment deletes a single comment by extract its ID
func (c *Client) DeleteComment(comment *reddit.Comment, logChan chan<- Log) error {
	ctx := context.Background()

	if c.dryRun {
		runtime.LogInfof(c.ctx, "[Reddit Client] Dry run: Comment would be deleted: %s", comment.FullID)
		logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Dry run: Comment would be deleted: %s", comment.FullID), Done: false, Error: nil}
		return nil
	}

	resp, err := c.client.Comment.Delete(ctx, comment.FullID)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	runtime.LogDebugf(c.ctx, "[Reddit Client] Comment deleted: %s", comment.FullID)
	logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Comment deleted: %s", comment.FullID), Done: false, Error: nil}

	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// EditComment edits a single comment, setting the content to the given words
func (c *Client) EditComment(comment *reddit.Comment, words string, logChan chan<- Log) error {
	ctx := context.Background()

	if c.dryRun {
		runtime.LogInfof(c.ctx, "[Reddit Client] Dry run: Comment would be edited: %s", comment.FullID)
		logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Dry run: Comment would be edited: %s", comment.FullID), Done: false, Error: nil}
		return nil
	}

	_, resp, err := c.client.Comment.Edit(ctx, comment.FullID, words)
	if err != nil {
		return err
	}

	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	secUntilReset := resp.Header.Get("X-Ratelimit-Reset")

	runtime.LogDebugf(c.ctx, "[Reddit Client] Comment edited: %s", comment.FullID)
	logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Comment edited: %s", comment.FullID), Done: false, Error: nil}

	// if we're out of requests, sleep
	err = evaluateLimit(remaining, secUntilReset)

	return err
}

// RetrieveAllUserComments retrieves all comments for the authenticated user
func (c *Client) RetrieveAllUserComments(logChan chan<- Log) ([]*reddit.Comment, error) {
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
			for _, comment := range pageComments {
				runtime.LogDebugf(c.ctx, "[Reddit Client] Comment fetched: %s", comment.FullID)
				logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Comment fetched: %s", comment.FullID), Done: false, Error: nil}
			}
		}

		comments = append(comments, pageComments...)

		runtime.LogDebugf(c.ctx, "[Reddit Client] Comments fetched: %d", len(comments))
		logChan <- Log{Time: time.Now(), Message: fmt.Sprintf("Fetched %d comments; total: %d", len(pageComments), len(comments)), Done: false, Error: nil}

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
