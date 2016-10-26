package gorpc

import (
	"net/http"
	"reflect"
	"fmt"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Secret   string
	services *serviceMap
	OnExecute func(action string, req interface{}, res interface{})
	OnError func(err error)
}

type Contexter interface {
	Method() string
	RequestUrl()string
	RequestHeader(key string) string
	RequestBody() ([]byte, error)
	ResponseHeader(key, value string)
	ResponseWrite(statusCode int, body []byte) error
}

func NewServer(secret string) *Server {
	return &Server{Secret: secret, services: new(serviceMap)}
}

func (this *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := NewHttpContext(writer, request)
	this.Execute(ctx)
}


func (this *Server) ServeFastHTTP(httpCtx *fasthttp.RequestCtx) {
	ctx := NewFastContext(httpCtx)
	this.Execute(ctx)
}

func (this *Server) Execute(ctx Contexter) {
	reply, err := this.paserAndExecute(ctx)
	if err != nil {
		if this.OnError != nil{
			this.OnError(fmt.Errorf("rpc execute error: %v", err))
		}
	}
	err = this.write(ctx, reply, err)
	if err != nil {
		if this.OnError != nil{
			this.OnError(fmt.Errorf("rpc write error: %v", err))
		}
	}
}

func (this *Server) Register(services ...interface{}) error {
	for _, service := range services {
		err := this.RegisterWithName(service, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Server) RegisterWithName(service interface{}, name string) error {
	return this.services.register(service, name)
}

func (this *Server) paserAndExecute(ctx Contexter) (interface{}, error) {
	if ctx.Method() != "POST" {
		return nil, ErrURLInvalid
	}
	sign := ctx.RequestHeader("sign")
	timestamp := ctx.RequestHeader("timestamp")
	action := ctx.RequestHeader("action")
	if sign == "" || timestamp == "" || action == "" {
		return nil, ErrURLInvalid
	}
	byteBody, err := ctx.RequestBody()
	if err != nil {
		return nil, err
	}
	if this.Secret != "" {
		reqSign := makeSign(timestamp + string(byteBody), this.Secret)
		if reqSign != sign {
			return nil, ErrPasswordIncorrect
		}
	}
	serviceSpec, methodSpec, errGet := this.services.get(action)
	if errGet != nil {
		return nil, fmt.Errorf("rpc found action %v error %v", action, errGet)
	}
	refArgs := reflect.New(methodSpec.argsType)
	args := refArgs.Interface()
	if len(byteBody) > 0 {
		err := json.Unmarshal(byteBody, &args)
		if err != nil {
			return nil, err
		}
	}

	// Call the service method.
	refReply := reflect.New(methodSpec.replyType)
	reply := refReply.Interface()

	// omit the HTTP request if the service method doesn't accept it
	var errValue []reflect.Value
	errValue = methodSpec.method.Func.Call([]reflect.Value{
		serviceSpec.rcvr,
		refArgs,
		refReply,
	})
	// Cast the result to error if needed.
	var errResult error
	errInter := errValue[0].Interface()
	if errInter != nil {
		errResult = errInter.(error)
	}
	if this.OnExecute != nil{
		this.OnExecute(action, args, reply)
	}
	return reply, errResult
}

func (this *Server) write(ctx Contexter, reply interface{}, err error) error{
	var body []byte
	if err == nil && reply != nil {
		body, err = json.Marshal(reply)
	}
	ctx.ResponseHeader("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		ctx.ResponseHeader("msg", err.Error())
	}
	return ctx.ResponseWrite(http.StatusOK, body)
}


