package gorpc

import (
	"errors"
	"log"
	"runtime"
)

var (
	ErrSignInvalid = errors.New("Sign invalid!")
	ErrURLInvalid = errors.New("URL error!")
	ErrMethodNotFound = errors.New("Method not found or not register!")
)



func logError(err error)  {
	log.Printf("error: %v", err)
	for i := 1; i < 10; i++{
		_, file, line, ok := runtime.Caller(i)
		if !ok{
			break
		}
		log.Printf("at file: %v, line: %v", file, line)
	}

}


func logMsg(errmsg string)  {
	log.Printf("error: %v", errmsg)
	for i := 1; i < 10; i++{
		_, file, line, ok := runtime.Caller(i)
		if !ok{
			break
		}
		log.Printf("at file: %v, line: %v", file, line)
	}

}