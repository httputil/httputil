package httputil

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	StatusCodes = []int{
		http.StatusContinue,
		http.StatusSwitchingProtocols,
		http.StatusProcessing,
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,
		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusUseProxy,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusProxyAuthRequired,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusLengthRequired,
		http.StatusPreconditionFailed,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestURITooLong,
		http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable,
		http.StatusExpectationFailed,
		http.StatusTeapot,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusLocked,
		http.StatusFailedDependency,
		http.StatusTooEarly,
		http.StatusUpgradeRequired,
		http.StatusPreconditionRequired,
		http.StatusTooManyRequests,
		http.StatusRequestHeaderFieldsTooLarge,
		http.StatusUnavailableForLegalReasons,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended,
		http.StatusNetworkAuthenticationRequired,
	}
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
			format = "plain"
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
	Message string `json:"message" xml:"message"`
}

func setHeaders(c echo.Context) {
	var (
		req    = c.Request()
		method = req.Method
		accept = req.Header.Get(echo.HeaderAccept)
		isTLS  = c.IsTLS()
		realIP = c.RealIP()
	)

	rh := c.Response().Header()
	rh.Set("X-Request-Method", method)
	rh.Set("X-Request-Accept", accept)
	rh.Set("X-Request-TLS", boolToHeaderValue(isTLS))
	rh.Set("X-Real-IP", realIP)
}

func boolToHeaderValue(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
