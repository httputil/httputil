package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/httputil/httputil"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	Debug = os.Getenv("RELEASE") != "1"
)

func main() {
	cfg := DefaultConfig
	if err := cfg.Load(); err != nil {
		log.Fatal("load configuration failed: ", err)
	}

	s := newServer()
	go func() {
		if err := s.Start(":" + cfg.PORT); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe exited: ", err)
		}
	}()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh
	log.Print("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Print("shutdown completed")
}

type Config struct {
	PORT string
}

func (c *Config) Load() error {
	if port := os.Getenv("PORT"); port != "" {
		c.PORT = port
	}
	return nil
}

var DefaultConfig = Config{PORT: "8080"}

func newServer() *echo.Echo {
	s := echo.New()
	s.HideBanner = !Debug
	s.HidePort = !Debug

	s.Use(middleware.CORS())

	for _, r := range httputil.DefaultRoutes {
		switch {
		case len(r.Methods) == 0:
			s.Any(r.Path, r.Handler)
		default:
			s.Match(r.Methods, r.Path, r.Handler)
		}
	}
	return s
}
