package zetka

import (
	"net/http"
	"net/http/pprof"
)

type Option func(*Client) error

func WithPprof(addr string) func(c *Client) error {
	//if _, err :=; err != nil {
	//	// TODO(tobbbles): is this dumb?
	//	return func(_ *Client) error {
	//		return errors.Wrap(err, "invalid url given")
	//	}
	//}

	return func(c *Client) error {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		c.srv = &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		return nil
	}
}
