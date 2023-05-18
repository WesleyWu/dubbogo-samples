package service

import (
	"context"

	api "dubbo-go-server/api"
	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/dubbogo/gost/log/logger"
)

type GreeterServerImpl struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterServerImpl) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	logger.Infof("Dubbo-go GreeterProvider get user name = %s\n", in.Name)
	return &api.HelloReply{Message: "Hello " + in.Name}, nil
}

func init() {
	config.SetProviderService(&GreeterServerImpl{})
}
