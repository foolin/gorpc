gorpc is a lightweight RPC over HTTP services, use password sign for security, providing access to the exported methods of an object through HTTP requests.

Features
--------
* Lightweight RPC
* Http service
* Password for security

Usage
---------

Install:

    go get github.com/foolin/gorpc


Example
---------

## Service

```go

package rpcvo

type UserService struct {}

func (this *UserService) Say(name *string, reply *string) error {
	*reply = "Hello, " + *name + "!"
	return nil
}


```


## Server
```go

package main
import (
	"net/http"
	"log"
	"github.com/foolin/gorpc"
	"github.com/foolin/gorpc/example/rpcvo"
)


func main()  {
	log.SetFlags(log.Lshortfile|log.LstdFlags)

	secrect := "1q2w3e" //password

	server := gorpc.NewServer(secrect)

	//=========== register rpc service =======//
	//Register batch
	err := server.Register(
		new(rpcvo.UserService),
	)
	if err != nil {
		log.Panicf("register error %v", err)
	}

	//=========== create http server =======//
	mux := http.NewServeMux()
	mux.Handle("/rpc", server)
	log.Println("listen: 5080")
	log.Print(http.ListenAndServe(":5080", mux))
}


```


## Client

```go

package main
import (
	"github.com/foolin/gorpc"
	"log"
	"github.com/foolin/gorpc/example/rpcvo"
)


func main()  {
	//log.SetFlags(log.Lshortfile|log.LstdFlags)

	client := gorpc.NewClient("http://127.0.0.1:5080/rpc", "1q2w3e")
	var err error

	//============= user service =============//
	var sayReply string
	err = client.Call("UserService.Say", "Foolin", &sayReply)
	if err != nil {
		log.Printf("UserService.Say error: %v", err)
	}else{
		log.Printf("UserService.Say result: %v", sayReply)
	}

}


```

## Result

    2016/04/14 10:25:58 UserService.Say result: Hello, Foolin!
