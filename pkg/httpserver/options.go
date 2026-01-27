package httpserver

import "time"

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Option configures Server through functional options pattern.
type Option func(*options)

// options holds server configuration parameters.
type options struct {
	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

// defaultOptions returns sane configuration defaults.
func defaultOptions() *options {
	return &options{
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}
}

// WithAddress sets server bind address.
func WithAddress(addr string) Option {
	return func(o *options) {
		o.address = addr
	}
}

// WithReadTimeout sets maximum duration for reading request headers/body.
func WithReadTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.readTimeout = timeout
	}
}

// WithWriteTimeout sets maximum duration for writing response.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.writeTimeout = timeout
	}
}

// WithShutdownTimeout sets graceful shutdown deadline.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.shutdownTimeout = timeout
	}
}
