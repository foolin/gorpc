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

	//============= calculator service =============//
	mutiArgs := &rpcvo.CalculatorArgs{
		A: 15,
		B: 100,
	}
	var mutiResult rpcvo.CalculatorReply
	err = client.Call("CalculatorService.Multiply", mutiArgs, &mutiResult)
	if err != nil {
		log.Printf("CalculatorService.Multiply error: %v", err)
	}else{
		log.Printf("CalculatorService.Multiply result: %v", mutiResult.Result)
	}

	//============= calculator service with name =============//
	divArgs := &rpcvo.CalculatorArgs{
		A: 15,
		B: 0,
	}
	var divResult rpcvo.CalculatorReply
	err = client.Call("Cal.Div", divArgs, &divResult)
	if err != nil {
		log.Printf("Cal.Div error: %v", err)
	}else{
		log.Printf("Cal.Div result: %v", divResult.Result)
	}
}
