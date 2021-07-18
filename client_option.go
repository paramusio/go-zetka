package zetka

type Option func(*Client) error

// WithPrometheus will attach a /metrics endpoint for Prometheus metrics on the DefaultPort
func WithPrometheus() Option {
	return func(c *Client) error {
		return registerMetrics()
	}
}
