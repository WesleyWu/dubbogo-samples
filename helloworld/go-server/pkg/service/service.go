package service

import (
	"context"
	"fmt"
	"os"

	api "dubbo-go-server/api"
	"dubbo.apache.org/dubbo-go/v3/config"
	"github.com/dubbogo/gost/log/logger"
)

var hostName string

type GreeterServerImpl struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterServerImpl) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	logger.Infof("Dubbo-go GreeterProvider get user name = %s\n", in.Name)
	return &api.HelloReply{Message: fmt.Sprintf("Hello %s, from %s", in.Name, hostName)}, nil
}

func init() {
	config.SetProviderService(&GreeterServerImpl{})
	var err error
	hostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}
