package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	greetv1 "github.com/sisisin-sandbox/try-connect-go/gen/greet/v1"
	"github.com/sisisin-sandbox/try-connect-go/gen/greet/v1/greetv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/proto"
)

type GreetServer struct{}

// GreetError implements greetv1connect.GreetServiceHandler.
func (s *GreetServer) GreetError(context.Context, *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	err := connect.NewError(connect.CodeUnavailable, fmt.Errorf("greet.v1.GreetService.GreetError is not implemented"))
	detail := &errdetails.DebugInfo{
		Detail: "this is debug info",
	}
	if detail, detailErr := connect.NewErrorDetail(detail); detailErr == nil {
		err.AddDetail(detail)
	}

	return nil, err
}

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
	interceptors := connect.WithInterceptors(NewValidateInterceptor())

	mux.Handle(greetv1connect.NewGreetServiceHandler(NewGreetServer(), interceptors))

	http.ListenAndServe(
		"localhost:8181",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func NewValidateInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			v, err := protovalidate.New()
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to initialize validator: %v", err))
			}
			msg, ok := req.Any().(proto.Message)
			if !ok {
				return nil, errors.New("failed to type assertion proto.Message")
			}
			if err = v.Validate(msg); err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("validation failed: %v", err))
			}
			return next(ctx, req)

		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
