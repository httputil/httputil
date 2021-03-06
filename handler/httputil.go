package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/dnscache"
)

var (
	statusCodes = []int{
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
		//http.StatusTeapot,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusLocked,
		http.StatusFailedDependency,
		//http.StatusTooEarly,
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

	DefaultRoutes = []Route{}
)

func init() {
	for _, code := range statusCodes {
		DefaultRoutes = append(DefaultRoutes,
			Route{
				Path:    strconv.Itoa(code),
				Handler: HandleHTTPStatus,
			},
			Route{
				Path:    strconv.Itoa(code) + "/:format",
				Handler: HandleHTTPStatus,
			},
		)
	}

	DefaultRoutes = append(DefaultRoutes,
		Route{
			Path:    "ip",
			Handler: HandleIP,
		},
		Route{
			Path:    "ip/:format",
			Handler: HandleIP,
		},
	)
}

type Route struct {
	Methods []string
	Path    string
	Handler echo.HandlerFunc
}

func NewServer(debug bool) *echo.Echo {
	s := echo.New()
	s.HideBanner = !debug
	s.HidePort = !debug

	s.Use(middleware.Recover())
	s.Use(middleware.CORS())
	s.Use(NoCacheHeaderMiddleware())
	s.Use(CachedResolverMiddleware(&dnscache.Resolver{Timeout: time.Second}))

	for _, r := range DefaultRoutes {
		switch {
		case len(r.Methods) == 0:
			s.Any(r.Path, r.Handler)
		default:
			s.Match(r.Methods, r.Path, r.Handler)
		}
	}
	return s
}

func setHeaders(c echo.Context) {
	var (
		req    = c.Request()
		method = req.Method
		accept = req.Header.Get(echo.HeaderAccept)
		realIP = c.RealIP()
	)

	rh := c.Response().Header()
	rh.Set("X-Request-Method", method)
	rh.Set("X-Request-Accept", accept)
	rh.Set(echo.HeaderXRealIP, realIP)
}

func boolToHeaderValue(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
