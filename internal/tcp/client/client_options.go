package client

import "time"

type ClientOption func(*Client)

func WithIdleTimeout(idleTimeout time.Duration) ClientOption {
	return func(c *Client) {
		c.idleTimeout = idleTimeout
	}
}

func WithBufferSize(bufferSize int) ClientOption {
	return func(c *Client) {
		c.bufferSize = bufferSize
	}
}
