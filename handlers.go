package httputil

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func HandleHTTPStatus(c echo.Context) error {
	setHeaders(c)

	paths := strings.Split(c.Path()[1:], "/")
	code, err := strconv.Atoi(paths[0])
	if err != nil {
		return err
	}

	switch code {
	case http.StatusNoContent, http.StatusNotModified:
		return c.NoContent(code)
	case
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		return c.Redirect(code, "https://httputil.dev/200")
	}

	// Create response
	status := http.StatusText(code)
	if status == "" {
		return errors.New("")
	}
	resp := HTTPStatusResponse{Message: status}

	// Send response
	format := c.Param("format")

	accept := c.Request().Header.Get(echo.HeaderAccept)
	if format == "" {
		switch {
		case strings.Contains(accept, echo.MIMEApplicationJSON):
			format = "json"
		case strings.Contains(accept, echo.MIMEApplicationXML), strings.Contains(accept, echo.MIMETextXML):
			format = "xml"
		default:
			format = "text"
		}
	}
	switch format {
	case "json":
		return c.JSON(code, resp)
	case "xml":
		return c.XML(code, resp)
	default:
		return c.String(code, resp.Message)
	}
}

type HTTPStatusResponse struct {
	Message string `json:"message" xml:"Message"`
}

func HandleIP(c echo.Context) error {
	setHeaders(c)
	ctx := c.Request().Context()
	rslv := ResolverFromEcho(c)

	var resp IPResponse
	resp.IP = c.RealIP()
	hosts, _ := rslv.LookupAddr(ctx, resp.IP)
	if len(hosts) > 0 {
		resp.Host = hosts[0]
	}

	return c.JSON(http.StatusOK, resp)
}

type IPResponse struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}
