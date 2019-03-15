package handler

import (
	"context"
	"net"

	"github.com/labstack/echo/v4"
)

const resolverCtxKey = "httputil.resolver"

type Resolver interface {
	LookupAddr(ctx context.Context, in string) (names []string, err error)
}

func ResolverFromEcho(ctx echo.Context) Resolver {
	r, ok := ctx.Get(resolverCtxKey).(Resolver)
	if !ok {
		return net.DefaultResolver
	}
	return r
}

func CachedResolverMiddleware(r Resolver) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(resolverCtxKey, r)
			return next(c)
		}
	}
}
