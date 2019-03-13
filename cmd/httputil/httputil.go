package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	s := echo.New()
	s.HideBanner = !Debug
	s.HidePort = !Debug

	s.Use(middleware.CORS())

	for _, code := range httputil.StatusCodes {
		s.Any(strconv.Itoa(code), httputil.HandleHTTPStatus)
		s.Any(strconv.Itoa(code)+"/:format", httputil.HandleHTTPStatus)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := s.Start(":" + port); err != nil && err != http.ErrServerClosed {
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

func wrapHandler(h http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
