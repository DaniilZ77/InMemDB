package client

import "time"

type ClientOption func(*Client)

func WithReadTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.readTimeout = timeout
	}
}

func WithBufferSize(bufferSize int) ClientOption {
	return func(c *Client) {
		c.bufferSize = bufferSize
	}
}
