package gorpc

import (
	"net/http"
	"io/ioutil"
	"fmt"
)


type HttpContext struct {
	request  *http.Request
	response http.ResponseWriter
}

func NewHttpContext(writer http.ResponseWriter, req *http.Request) *HttpContext {
	return &HttpContext{request: req, response: writer}
}

func (ctx *HttpContext) Method()  string{
	return ctx.request.Method
}


func (ctx *HttpContext) RequestUrl()  string{
	return ctx.request.RequestURI
}

func (ctx *HttpContext) RequestHeader(key string) string{
	return ctx.request.Header.Get(key)
}


func (ctx *HttpContext) RequestBody() ([]byte, error){
	return ioutil.ReadAll(ctx.request.Body)
}

func (ctx *HttpContext) ResponseHeader(key, value string){
	ctx.response.Header().Set(key, value)
}

func (ctx *HttpContext) ResponseWrite(statusCode int, body []byte) error{
	//header必须在后面
	ctx.response.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(ctx.response, body)
	return err
}