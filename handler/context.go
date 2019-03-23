package handler

import "github.com/labstack/echo/v4"
import "github.com/ugorji/go/codec"

var (
	msgpackHandle = &codec.MsgpackHandle{}
)

// Context is wrapped echo.Context
type Context interface {
	echo.Context
}

type ctx struct {
	echo.Context
}

func Wrap(c echo.Context) *ctx {
	return &ctx{Context: c}
}

// Msgpack sends an MessagePack response with status code.
func (c *ctx) Msgpack(code int, i interface{}) (err error) {
	resp := c.Response()
	resp.Header().Set(echo.HeaderContentType, echo.MIMEApplicationMsgpack)
	resp.WriteHeader(code)
	return codec.NewEncoder(resp, msgpackHandle).Encode(i)
}
