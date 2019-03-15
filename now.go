package functions

import (
	"net/http"
	"os"

	"github.com/httputil/httputil/handler"
)

var (
	defaultServer = handler.NewServer(os.Getenv("RELEASE") != "1")
)

func HandleAll(w http.ResponseWriter, r *http.Request) {
	defaultServer.ServeHTTP(w, r)
}
