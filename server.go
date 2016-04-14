package gorpc

import (
	"net/http"
	"io/ioutil"
	"reflect"
	"fmt"
	"encoding/json"
)

type Server struct {
	Secret string
	services   *serviceMap
}


func NewServer(secret string) *Server {
	return &Server{Secret: secret, services: new(serviceMap)}
}

func (this *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	reply, err := this.execute(request)
	this.write(writer, reply, err)
}

func (this *Server) Register(services ...interface{}) error  {
	for _, service := range services{
		err := this.RegisterWithName(service, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Server) RegisterWithName(service interface{}, name string) error  {
	return this.services.register(service, name)
}


func (this *Server) execute(request *http.Request) (interface{}, error)  {
	if request.Method != "POST" {
		return nil, ErrURLInvalid
	}
	sign := request.Header.Get("sign")
	timestamp := request.Header.Get("timestamp")
	action := request.Header.Get("action")
	if sign == "" || timestamp == "" || action == ""{
		return nil, ErrURLInvalid
	}
	byteBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	reqSign := makeSign(timestamp + string(byteBody), this.Secret)
	if reqSign != sign{
		return nil, ErrSignInvalid
	}
	serviceSpec, methodSpec, errGet := this.services.get(action)
	if errGet != nil {
		logError(errGet)
		return nil, ErrMethodNotFound
	}
	refArgs := reflect.New(methodSpec.argsType)
	args := refArgs.Interface()
	if len(byteBody) > 0{
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


func (this *Server) write(w http.ResponseWriter, reply interface{}, err error)  {
	body := ""
	if err == nil && reply != nil{
		var bytes []byte
		bytes, err = json.Marshal(reply)
		if err != nil {
			logError(err)
		}else{
			body = string(bytes)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil{
		w.Header().Set("msg", err.Error())
	}
	//header必须在后面
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}


