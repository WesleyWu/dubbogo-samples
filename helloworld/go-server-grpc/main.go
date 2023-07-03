package main

import (
	"context"

	"github.com/WesleyWu/go-lifespan/lifespan"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/os/gctx"
	_ "go-server-grpc/service"
)

func main() {
	s := grpcx.Server.New()
	ctx := context.WithValue(gctx.GetInitCtx(), "grpc.Server", s.Server)
	lifespan.OnBootstrap(ctx)
	s.Run()
	lifespan.OnShutdown(ctx)
}
