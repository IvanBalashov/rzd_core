package http

import (
	"net"
	"net/http"
	"time"
)

func NewServer(handler http.Handler, host, port string) *http.Server {
	server := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 10,
	}

	return server
}
