/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"sync"
	"time"

	api "dubbo-go-client/api"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	_ "dubbo.apache.org/dubbo-go/v3/registry/xds"
	"github.com/dubbogo/gost/log/logger"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

var grpcGreeterImpl = new(api.GreeterClientImpl)

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-client/conf/dubbogo.yml
func main() {
	config.SetConsumerService(grpcGreeterImpl)
	if err := config.Load(); err != nil {
		panic(err)
	}

	SayHelloTo("Wesley Wu")
	s := g.Server()
	s.BindHandler("/hello", SayHelloHandler)
	s.BindHandler("/benchmark", BenchmarkHandler)
	s.Run()
}

func SayHelloTo(name string) {
	req := &api.HelloRequest{
		Name: name,
	}
	reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("client response result: %v\n", reply)
}

func SayHelloHandler(r *ghttp.Request) {
	name := r.Get("name", "Wesley Wu").String()
	req := &api.HelloRequest{
		Name: name,
	}
	reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
	if err != nil {
		logger.Error(err)
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
			reply, err := grpcGreeterImpl.SayHello(context.Background(), req)
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
	r.Response.WritefExit("call rpc %d times, %s elapsed milliseconds", times, elapsed)
}
