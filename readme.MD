# DubboGo Samples 

## helloworld
### go-server
dubbogo 微服务 provider，用Go语言实现，支持
- 本地运行，服务注册到 etcdv3，可以换成 nacos/polaris 等
- 部署到 k8s istio 服务网格，支持多个副本同时提供服务

#### 接口描述：
- SayHello: 根据传入参数 name，输出 "Hello " + name 

#### 安装部署
参考文档：https://cn.dubbo.apache.org/zh-cn/overview/mannual/golang-sdk/tutorial/governance/service-mesh/deploy
```shell
make build
make deploy
# 如果有修改，需要升级版本，请调用 make upgrade
```

### go-client: 
dubbogo 微服务 consumer & HTTP server (基于 GoFrame)
用Go语言实现，并启动一个 HTTP Web Server，支持
- 本地运行
- 部署到 k8s

功能：
- 通过 k8s istio / etcdv3(本地运行模式) 自动发现 Provider，多 Provider 副本由 k8s 进行负载均衡
- 提供 /hello http 端点，接收带 ?name=xxx 查询的 GET 请求，调用 provider 的 SayHello 方法，并返回结果 
#### 安装部署
参考文档：https://cn.dubbo.apache.org/zh-cn/overview/mannual/golang-sdk/tutorial/governance/service-mesh/deploy

```shell
make build
make deploy
# 如果有修改，需要升级版本，请调用 make upgrade
```

### java-client
dubbogo 微服务 client & HTTP server (基于SprintBoot)

### python-client
dubbogo 微服务 client & HTTP server (基于Flask)
