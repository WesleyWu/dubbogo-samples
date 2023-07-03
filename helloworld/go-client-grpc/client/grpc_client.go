package client

import (
	"context"
	"time"

	"github.com/WesleyWu/go-lifespan/lifespan"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/genv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Conn *grpc.ClientConn
	err  error
)

func init() {
	lifespan.RegisterBootstrapHook("grpc", true, func(ctx context.Context) error {
		Conn, err = makeConnection(ctx)
		return err
	})
	lifespan.RegisterShutdownHook("grpc", true, func(ctx context.Context) error {
		_ = Conn.Close()
		return nil
	})
}

func makeConnection(ctx context.Context) (conn *grpc.ClientConn, err error) {
	serviceEndpoint := genv.Get("SERVICE_ENDPOINT", "").String()
	if serviceEndpoint == "" {
		err = gerror.New("SERVICE_ENDPOINT env not set")
		return
	}
	timeoutSeconds := genv.Get("CONNECT_TIMEOUT", 5).Int()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	conn, err = grpc.DialContext(timeoutCtx, serviceEndpoint,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, gerror.Newf("grpc connection to (%s) failed: %+v", serviceEndpoint, err)
	}
	return
}
