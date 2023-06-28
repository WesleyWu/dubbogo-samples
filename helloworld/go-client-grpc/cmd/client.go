/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Binary main implements a client for Greeter service using gRPC's client-side
// support for xDS APIs.
package main

import (
	"context"
	"sync"
	"time"

	api "dubbo-go-client/api"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/xds" // To install the xds resolvers and balancers.
)

var (
	conn    *grpc.ClientConn
	greeter api.GreeterClient
)

func makeConnection(ctx context.Context) (conn *grpc.ClientConn, grpcGreeterImpl api.GreeterClient, err error) {
	serviceEndpoint := genv.Get("SERVICE_ENDPOINT", "").String()
	if serviceEndpoint == "" {
		err = gerror.New("SERVICE_ENDPOINT env not set")
		return
	}
	conn, err = grpc.Dial(serviceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		g.Log().Fatalf(ctx, "grpc.Dial(%s) failed: %v", serviceEndpoint, err)
		return
	}
	grpcGreeterImpl = api.NewGreeterClient(conn)
	return
}

func main() {
	var err error
	ctx := gctx.New()
	conn, greeter, err = makeConnection(ctx)
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	SayHelloTo(ctx, "Wesley Wu")
	s := g.Server()
	s.BindHandler("/hello", SayHelloHandler)
	s.BindHandler("/benchmark", BenchmarkHandler)
	s.Run()
}

func SayHelloTo(ctx context.Context, name string) {
	req := &api.HelloRequest{
		Name: name,
	}
	reply, err := greeter.SayHello(context.Background(), req)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	g.Log().Infof(ctx, "client response result: %v\n", reply)
}

func SayHelloHandler(r *ghttp.Request) {
	name := r.Get("name", "Wesley Wu").String()
	req := &api.HelloRequest{
		Name: name,
	}
	reply, err := greeter.SayHello(context.Background(), req)
	if err != nil {
		g.Log().Error(r.Context(), err)
		r.Response.WriteStatusExit(500, err.Error())
		return
	}
	r.Response.WritelnExit(reply.Message)
}

func BenchmarkHandler(r *ghttp.Request) {
	name := r.Get("name", "Wesley Wu").String()
	req := &api.HelloRequest{
		Name: name,
	}
	times := r.Get("times", 100).Int()

	var wg sync.WaitGroup
	start := gtime.Now()
	mu := sync.Mutex{}
	for index := 0; index < times; index++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			reply, err := greeter.SayHello(context.Background(), req)
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
				r.Response.WriteStatusExit(500, err.Error())
				return
			}
			mu.Lock()
			r.Response.Writeln(reply.Message)
			mu.Unlock()
		}(r.Context())
	}
	wg.Wait()
	end := gtime.Now()
	elapsed := int64(end.Sub(start) / time.Millisecond)
	r.Response.WritefExit("call rpc %d times, %d elapsed milliseconds", times, elapsed)
}
