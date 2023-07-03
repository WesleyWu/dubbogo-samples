package service

import (
	"context"
	"fmt"
	"os"

	"github.com/WesleyWu/go-lifespan/lifespan"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	api "go-server-grpc/api"
	"google.golang.org/grpc"
)

var hostName string

func init() {
	lifespan.RegisterBootstrapHook("greeter", true, func(ctx context.Context) error {
		grpcServer := ctx.Value("grpc.Server").(*grpc.Server)
		api.RegisterGreeterServer(grpcServer, &GreeterServerImpl{})
		return nil
	})
}

type GreeterServerImpl struct {
	api.UnimplementedGreeterServer
}

func Register(s *grpcx.GrpcServer) {
}

func (s *GreeterServerImpl) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	//g.Log().Infof(ctx, "Dubbo-go GreeterProvider get user name = %s\n", in.Name)
	return &api.HelloReply{Message: fmt.Sprintf("Hello %s, from %s", in.Name, hostName)}, nil
}

func init() {
	var err error
	hostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}
