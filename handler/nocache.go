package handler

import "github.com/labstack/echo/v4"

var (
	noCacheReqHeaderToDel = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}
	noCacheResHeaderToSet = map[string]string{
		"Expires":       "Thu, 01 Jan 1970 09:00:00 JST",
		"Cache-Control": "no-cache, private, max-age=0",
		"Pragma":        "no-cache",
	}
)

func NoCacheHeaderMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			for _, h := range noCacheReqHeaderToDel {
				req.Header.Del(h)
			}

			res := c.Response()
			for k, v := range noCacheResHeaderToSet {
				res.Header().Set(k, v)
			}

			return next(c)
		}
	}
}
