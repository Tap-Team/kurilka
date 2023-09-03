package app

import (
	"net/http"
	"time"

	"github.com/Tap-Team/kurilka/internal/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Server(handler http.Handler, cnf config.ServerConfig) *http.Server {
	h2s := &http2.Server{
		IdleTimeout: 10 * time.Second,
	}
	s := http.Server{
		Addr:    cnf.Addr(),
		Handler: h2c.NewHandler(handler, h2s),
	}
	return &s
}
