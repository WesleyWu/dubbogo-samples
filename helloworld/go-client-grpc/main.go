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
	"sync/atomic"
	"time"

	api "dubbo-go-client/api"
	"dubbo-go-client/greeter"
	"github.com/WesleyWu/go-lifespan/lifespan"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/panjf2000/ants/v2"
	_ "google.golang.org/grpc/xds" // To install the xds resolvers and balancers.
)

func main() {
	ctx := gctx.New()
	lifespan.OnBootstrap(ctx)
	SayHelloTo(ctx, "Wesley Wu")
	s := g.Server()
	s.BindHandler("/hello", SayHelloHandler)
	s.BindHandler("/benchmark", BenchmarkHandler)
	s.Run()
	lifespan.OnShutdown(ctx)
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
	n := r.Get("n", 10).Int()
	defer ants.Release()
	req := &api.HelloRequest{
		Name: name,
	}
	times := r.Get("times", 100).Int()

	var wg sync.WaitGroup
	start := gtime.Now()

	ctx := r.Context()
	calls := atomic.Int64{}
	pool, err := ants.NewPoolWithFunc(n, func(i interface{}) {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			_, err := greeter.SayHello(ctx, i.(*api.HelloRequest))
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
				r.Response.WriteStatusExit(500, err.Error())
				return
			}
			calls.Add(1)
		}
	})
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		r.Response.WriteStatusExit(500, err.Error())
		return
	}
	defer pool.Release()
	for index := 0; index < times; index++ {
		wg.Add(1)
		_ = pool.Invoke(req)
	}
	wg.Wait()
	end := gtime.Now()
	elapsed := int64(end.Sub(start) / time.Millisecond)
	cps := int64(times) * 1000 / elapsed
	g.Log().Infof(ctx, "call rpc %d times, %d elapsed milliseconds, %d calls per second", calls.Load(), elapsed, cps)
	r.Response.WritefExit("call rpc %d times, %d elapsed milliseconds, %d calls per second", calls.Load(), elapsed, cps)
}
