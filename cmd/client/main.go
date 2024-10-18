package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	greetv1 "github.com/sisisin-sandbox/try-connect-go/gen/greet/v1"
	"github.com/sisisin-sandbox/try-connect-go/gen/greet/v1/greetv1connect"
)

func main() {

	client := greetv1connect.NewGreetServiceClient(http.DefaultClient, "http://localhost:8181", connect.WithGRPC())
	res, err := client.Greet(context.Background(), &connect.Request[greetv1.GreetRequest]{
		Msg: &greetv1.GreetRequest{Name: "greed"},
	})
	if err != nil {
		panic(err)
	}
	println(res.Msg.Greeting)

}
