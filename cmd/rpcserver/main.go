package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"flag"

	"github.com/alexandreroba/go-microservices/cmd/rpcserver/contract"
)

type HelloWorldHandler struct{}

func (h *HelloWorldHandler) HelloWorld(args *contract.HelloWorldRequest, reply *contract.HelloWorldResponse) error {
	reply.Message = "Hello" + args.Name
	return nil
}

func StartServer(port int) {
	helloWorld := &HelloWorldHandler{}
	rpc.Register(helloWorld)
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to listen on given port:%s", err))
	}
	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}

func CreateClient(port int) *rpc.Client {
	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return client
}

func PerformRequest(client *rpc.Client) contract.HelloWorldResponse {
	args := &contract.HelloWorldRequest{Name: "world"}
	var reply contract.HelloWorldResponse
	err := client.Call("HelloWorldHandler.HelloWorld", args, &reply)
	if err != nil {
		log.Fatal("error:", err)
	}
	return reply
}

var isServer = flag.Bool("server", false, "use the flag to run the command as server")
var port = flag.Int("port", 3000, "The rpc port")

func main() {
	var response contract.HelloWorldResponse
	if *isServer {
		StartServer(*port)
	} else {
		client := CreateClient(*port)
		response = PerformRequest(client)
		log.Println(response)
	}
}
