package zetka

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
)

// Client
type Client struct {
	token string

	gwurl    string
	compress bool

	srv *http.Server

	errs     chan error
	sequence atomic.Value
}

// New Client
func New(token string, opts ...Option) (*Client, error) {
	c := &Client{
		compress: true,

		errs: make(chan error),
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

	if c.srv != nil {
		go func(errc chan error) {
			if err := c.srv.ListenAndServe(); err != nil {
				errc <- fmt.Errorf("fatal error in internal server: %e", err)
				return
			}
		}(errc)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errc:
			return err
		}
	}
}
