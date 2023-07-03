package greeter

import (
	"context"

	api "dubbo-go-client/api"
	"dubbo-go-client/client"
	"github.com/WesleyWu/go-lifespan/lifespan"
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc"
)

var (
	greeter api.GreeterClient
)

func init() {
	lifespan.RegisterBootstrapHook("greeter", true, func(ctx context.Context) error {
		greeter = api.NewGreeterClient(client.Conn)
		return nil
	}, "grpc")
}

func SayHello(ctx context.Context, in *api.HelloRequest, opts ...grpc.CallOption) (*api.HelloReply, error) {
	if greeter == nil {
		return nil, gerror.New("No grpc client connection")
	}
	return greeter.SayHello(ctx, in, opts...)
}
