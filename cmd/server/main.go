package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	greetv1 "github.com/sisisin-sandbox/try-connect-go/gen/greet/v1"
	"github.com/sisisin-sandbox/try-connect-go/gen/greet/v1/greetv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}
func NewGreetServer() greetv1connect.GreetServiceHandler {
	return &GreetServer{}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle(greetv1connect.NewGreetServiceHandler(NewGreetServer()))

	http.ListenAndServe(
		"localhost:8181",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
