package gorpc

import (
	"net/http"
	"reflect"
	"fmt"
	"encoding/json"
	"log"
	"os"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Secret   string
	services *serviceMap
	Log      *log.Logger
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
	return &Server{Secret: secret, services: new(serviceMap), Log: log.New(os.Stdout, "rpc", log.LstdFlags)}
}

func (this *Server) SetLog(log *log.Logger) {
	this.Log = log
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
	this.write(ctx, reply, err)
	if err != nil {
		this.Log.Fatalf("rpc write error: %v", err)
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
	if this.Log != nil {
		this.Log.Printf("rpc listen url:%v, action:%v, sign:%v, timestamp:%v", ctx.RequestUrl(), action, sign, timestamp)
	}
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
		logError(errGet)
		return nil, ErrMethodNotFound
	}
	refArgs := reflect.New(methodSpec.argsType)
	args := refArgs.Interface()
	if len(byteBody) > 0 {
		err := json.Unmarshal(byteBody, &args)
		if err != nil {
			logMsg(fmt.Sprintf("action: %v, %v", action, err))
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

	return reply, errResult
}

func (this *Server) write(ctx Contexter, reply interface{}, err error){
	var body []byte
	if err == nil && reply != nil {
		body, err = json.Marshal(reply)
		if err != nil {
			logError(err)
		}
	}
	ctx.ResponseHeader("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		ctx.ResponseHeader("msg", err.Error())
	}
	err = ctx.ResponseWrite(http.StatusOK, body)
}


