package zetka

import (
	"context"
	"errors"
	"sync/atomic"
)

const DefaultPort = 4008

// Client
type Client struct {
	token string

	gwurl    string
	compress bool

	errs      chan error
	restartch chan struct{}

	sequence atomic.Value
}

// New Client
func New(token string, opts ...Option) (*Client, error) {
	c := &Client{
		compress: true,

		errs:      make(chan error),
		restartch: make(chan struct{}),
	}

	c.sequence.Store(int64(0))

	if len(token) == 0 {
		return nil, errors.New("no token provided to zetka")
	}
	c.token = token

	gwurl, err := gatewayURI(BaseURL, token, "6", JSON)
	if err != nil {
		return nil, err
	}
	c.gwurl = gwurl

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Start will block and listen to events and pass them to the results chan. Start will return at the first error
// received.
func (c *Client) Start(ctx context.Context, results chan *GatewayEvent) error {
	errc := make(chan error)

	go func(errc chan error, results chan *GatewayEvent) {
		if err := c.Receive(c.gwurl, results); err != nil {
			errc <- err
		}
	}(errc, results)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errc:
			return err
		}
	}
}
