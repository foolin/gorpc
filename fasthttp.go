package gorpc

import (
	"github.com/valyala/fasthttp"
	"net/http"
)


type FastContext struct {
	*fasthttp.RequestCtx
}

func NewFastContext(ctx *fasthttp.RequestCtx) *FastContext {
	return &FastContext{RequestCtx: ctx}
}

func (ctx *FastContext) Method()  string{
	return string(ctx.Request.Header.Method())
}


func (ctx *FastContext) RequestUrl()  string{
	return string(ctx.RequestURI())
}

func (ctx *FastContext) RequestHeader(key string) string{
	return string(ctx.Request.Header.Peek(key))
}


func (ctx *FastContext) RequestBody() ([]byte, error){
	return ctx.Request.Body(), nil
}

func (ctx *FastContext) ResponseHeader(key, value string){
	ctx.Response.Header.Set(key, value)
}

func (ctx *FastContext) ResponseWrite(statusCode int, body []byte) error{
	//header必须在后面
	ctx.SetStatusCode(http.StatusOK)
	_, err := ctx.Write(body)
	return err
}

