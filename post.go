package geddit

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// PostService handles communication with the link (post)
// related methods of the Reddit API.
type PostService interface {
	SubmitSelf(ctx context.Context, opts SubmitSelfOptions) (*Submitted, *Response, error)
	SubmitURL(ctx context.Context, opts SubmitURLOptions) (*Submitted, *Response, error)

	EnableReplies(ctx context.Context, id string) (*Response, error)
	DisableReplies(ctx context.Context, id string) (*Response, error)

	MarkNSFW(ctx context.Context, id string) (*Response, error)
	UnmarkNSFW(ctx context.Context, id string) (*Response, error)

	Spoiler(ctx context.Context, id string) (*Response, error)
	Unspoiler(ctx context.Context, id string) (*Response, error)

	Hide(ctx context.Context, ids ...string) (*Response, error)
	Unhide(ctx context.Context, ids ...string) (*Response, error)
}

// PostServiceOp implements the PostService interface.
type PostServiceOp struct {
	client *Client
}

var _ PostService = &PostServiceOp{}

type submittedLinkRoot struct {
	JSON struct {
		Data *Submitted `json:"data,omitempty"`
	} `json:"json"`
}

// Submitted is a newly submitted post on Reddit.
type Submitted struct {
	ID     string `json:"id,omitempty"`
	FullID string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
}

// SubmitSelfOptions are options used for selftext posts.
type SubmitSelfOptions struct {
	Subreddit string `url:"sr,omitempty"`
	Title     string `url:"title,omitempty"`
	Text      string `url:"text,omitempty"`

	FlairID   string `url:"flair_id,omitempty"`
	FlairText string `url:"flair_text,omitempty"`

	SendReplies *bool `url:"sendreplies,omitempty"`
	NSFW        bool  `url:"nsfw,omitempty"`
	Spoiler     bool  `url:"spoiler,omitempty"`
}

// SubmitURLOptions are options used for link posts.
type SubmitURLOptions struct {
	Subreddit string `url:"sr,omitempty"`
	Title     string `url:"title,omitempty"`
	URL       string `url:"url,omitempty"`

	FlairID   string `url:"flair_id,omitempty"`
	FlairText string `url:"flair_text,omitempty"`

	SendReplies *bool `url:"sendreplies,omitempty"`
	Resubmit    bool  `url:"resubmit,omitempty"`
	NSFW        bool  `url:"nsfw,omitempty"`
	Spoiler     bool  `url:"spoiler,omitempty"`
}

// SubmitSelf submits a self text post.
func (s *PostServiceOp) SubmitSelf(ctx context.Context, opts SubmitSelfOptions) (*Submitted, *Response, error) {
	type submit struct {
		SubmitSelfOptions
		Kind string `url:"kind,omitempty"`
	}
	return s.submit(ctx, &submit{opts, "self"})
}

// SubmitURL submits a link post.
func (s *PostServiceOp) SubmitURL(ctx context.Context, opts SubmitURLOptions) (*Submitted, *Response, error) {
	type submit struct {
		SubmitURLOptions
		Kind string `url:"kind,omitempty"`
	}
	return s.submit(ctx, &submit{opts, "link"})
}

func (s *PostServiceOp) submit(ctx context.Context, v interface{}) (*Submitted, *Response, error) {
	path := "api/submit"

	form, err := query.Values(v)
	if err != nil {
		return nil, nil, err
	}
	form.Set("api_type", "json")

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(submittedLinkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.JSON.Data, resp, nil
}

// EnableReplies enables inbox replies for a thing created by the client.
func (s *PostServiceOp) EnableReplies(ctx context.Context, id string) (*Response, error) {
	path := "api/sendreplies"

	form := url.Values{}
	form.Set("id", id)
	form.Set("state", "true")

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DisableReplies dsables inbox replies for a thing created by the client.
func (s *PostServiceOp) DisableReplies(ctx context.Context, id string) (*Response, error) {
	path := "api/sendreplies"

	form := url.Values{}
	form.Set("id", id)
	form.Set("state", "false")

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// MarkNSFW marks a post as NSFW.
func (s *PostServiceOp) MarkNSFW(ctx context.Context, id string) (*Response, error) {
	path := "api/marknsfw"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UnmarkNSFW unmarks a post as NSFW.
func (s *PostServiceOp) UnmarkNSFW(ctx context.Context, id string) (*Response, error) {
	path := "api/unmarknsfw"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Spoiler marks a post as a spoiler.
func (s *PostServiceOp) Spoiler(ctx context.Context, id string) (*Response, error) {
	path := "api/spoiler"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unspoiler unmarks a post as a spoiler.
func (s *PostServiceOp) Unspoiler(ctx context.Context, id string) (*Response, error) {
	path := "api/unspoiler"

	form := url.Values{}
	form.Set("id", id)

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Hide hides links with the specified ids.
func (s *PostServiceOp) Hide(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/hide"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Unhide unhides links with the specified ids.
func (s *PostServiceOp) Unhide(ctx context.Context, ids ...string) (*Response, error) {
	if len(ids) == 0 {
		return nil, errors.New("must provide at least 1 id")
	}

	path := "api/unhide"

	form := url.Values{}
	form.Set("id", strings.Join(ids, ","))

	req, err := s.client.NewPostForm(path, form)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}