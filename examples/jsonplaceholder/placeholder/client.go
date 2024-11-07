package placeholder

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zrbecker/fetch"
)

type Client struct {
	f *fetch.Fetcher
}

func NewClient(baseURL string, opts ...fetch.FetcherOption) *Client {
	c := &Client{
		f: fetch.NewFetcher(baseURL, opts...),
	}
	return c
}

func (c *Client) Posts(ctx context.Context, opts ...fetch.FetchOption) ([]Post, error) {
	var posts []Post
	if err := c.f.Fetch(ctx, "/posts", &posts, opts...); err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *Client) Post(ctx context.Context, id int, opts ...fetch.FetchOption) (Post, error) {
	var post Post
	if err := c.f.Fetch(ctx, fmt.Sprintf("/posts/%d", id), &post, opts...); err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c *Client) Comments(ctx context.Context, postID int, opts ...fetch.FetchOption) ([]Comment, error) {
	opts = append(
		[]fetch.FetchOption{
			fetch.Params(fetch.Param{Key: "postId", Value: strconv.Itoa(postID)}),
		},
		opts...,
	)
	var comments []Comment
	if err := c.f.Fetch(ctx, "/comments", &comments, opts...); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) CreatePost(ctx context.Context, params CreatePostParams, opts ...fetch.FetchOption) (Post, error) {
	opts = append(
		[]fetch.FetchOption{
			fetch.Method(http.MethodPost),
			fetch.Body(params),
			fetch.Headers(fetch.Header{Key: "Content-Type", Value: "application/json; charset=UTF-8"}),
		},
		opts...,
	)
	var post Post
	if err := c.f.Fetch(ctx, "/posts", &post, opts...); err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c *Client) UpdatePost(ctx context.Context, id int, params UpdatePostParams, opts ...fetch.FetchOption) (Post, error) {
	opts = append(
		[]fetch.FetchOption{
			fetch.Method(http.MethodPut),
			fetch.Body(params),
			fetch.Headers(fetch.Header{Key: "Content-Type", Value: "application/json; charset=UTF-8"}),
		},
		opts...,
	)
	var post Post
	if err := c.f.Fetch(ctx, fmt.Sprintf("/posts/%d", id), &post, opts...); err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c *Client) PatchPost(ctx context.Context, id int, params PatchPostParams, opts ...fetch.FetchOption) (Post, error) {
	opts = append(
		[]fetch.FetchOption{
			fetch.Method(http.MethodPatch),
			fetch.Body(params),
			fetch.Headers(fetch.Header{Key: "Content-Type", Value: "application/json; charset=UTF-8"}),
		},
		opts...,
	)
	var post Post
	if err := c.f.Fetch(ctx, fmt.Sprintf("/posts/%d", id), &post, opts...); err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c *Client) DeletePost(ctx context.Context, id int, opts ...fetch.FetchOption) error {
	opts = append(
		[]fetch.FetchOption{
			fetch.Method(http.MethodDelete),
		},
		opts...,
	)
	if err := c.f.Fetch(ctx, fmt.Sprintf("/posts/%d", id), nil, opts...); err != nil {
		return err
	}
	return nil
}
