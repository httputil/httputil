// +heroku install ./cmd/httputil/...
// +heroku goVersion go1.12

module github.com/httputil/httputil

go 1.12

require (
	github.com/labstack/echo/v4 v4.0.0
	github.com/rs/dnscache v0.0.0-20190225195841-509b4d5d3b47
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 // indirect
)
