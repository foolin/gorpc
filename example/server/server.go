package main
import (
	"net/http"
	"log"
	"xinyi.io/gameshared/module/gorpc"
	"xinyi.io/gameshared/module/gorpc/example/rpcvo"
)


func main()  {
	log.SetFlags(log.Lshortfile|log.LstdFlags)

	secrect := "1q2w3e" //password

	server := gorpc.NewServer(secrect)

	//=========== register rpc service =======//
	//Register batch
	err := server.Register(
		new(rpcvo.UserService),
		new(rpcvo.CalculatorService),
	)
	if err != nil {
		log.Panicf("register error %v", err)
	}
	//Register with name
	err = server.RegisterWithName(new(rpcvo.CalculatorService), "Cal")
	if err != nil {
		log.Panicf("register error %v", err)
	}

	//=========== create http server =======//
	mux := http.NewServeMux()
	mux.Handle("/rpc", server)
	/*
	//other way handler
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, req *http.Request) {

		//befor handler, write code here!!!

		server.ServeHTTP(w, req)

		//after handler, write code here!!!
	})
	*/
	log.Println("listen: 5080")
	log.Print(http.ListenAndServe(":5080", mux))
}


