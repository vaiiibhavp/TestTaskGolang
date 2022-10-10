package server

import (
	"net"
	"time"
)

// Option -.
type Option func(*HttpServer)

// Port -.
func Port(port string) Option {
	return func(s *HttpServer) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *HttpServer) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *HttpServer) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *HttpServer) {
		s.shutdownTimeout = timeout
	}
}
