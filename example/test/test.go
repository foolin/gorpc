package main

import (
	"reflect"
	"log"
	"unicode"
	"unicode/utf8"
)

func main() {
	log.SetFlags(log.LstdFlags)
	register(new(HelloService))
}



type Hello struct {
}

func (this *Hello) Sayer(args int, reply int) error {
	return nil
}

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct {}

func (h *HelloService) Say(args *HelloArgs, reply *HelloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}

func (h *HelloService) Hello(reply *HelloReply) error {
	reply.Message = "Hello" + "!"
	return nil
}


func (h *HelloService) mimi(reply *HelloReply) error {
	reply.Message = "Hello" + "!"
	return nil
}

var typeOfError   = reflect.TypeOf((*error)(nil)).Elem()

func register(receiver interface{})  {
	rcvr :=    reflect.ValueOf(receiver)
	rcvrType := reflect.TypeOf(receiver)

	name := reflect.Indirect(rcvr).Type().Name()
	log.Printf("object name: %v", name)
	for i := 0; i < rcvrType.NumMethod(); i++ {
		log.Printf("--------------")
		method := rcvrType.Method(i)
		mtype := method.Type
		log.Printf("method: %v, numIn: %v, pkpath: %v", method.Name, mtype.NumIn(), method.PkgPath)
		for j := 0; j < mtype.NumIn(); j++{
			arg := mtype.In(j)
			log.Printf("arg: %v", arg)
		}
		if method.Name == "Say"{
			reply := &HelloReply{}
			var errValue []reflect.Value
			errValue = method.Func.Call([]reflect.Value{
				rcvr,
				&HelloArgs{Who:"Foolin"},
				&reply,
			})
			// Cast the result to error if needed.
			var errResult error
			errInter := errValue[0].Interface()
			if errInter != nil {
				errResult = errInter.(error)
			}
			if errResult != nil {
				log.Printf("call error: %v", errResult)
			}else{
				log.Printf("call result: %v", reply)
			}
		}
	}
}

// isExported returns true of a string is an exported (upper case) name.
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// isExportedOrBuiltin returns true if a type is exported or a builtin.
func isExportedOrBuiltin(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}