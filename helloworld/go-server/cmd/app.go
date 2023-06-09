package main

import (
	_ "dubbo-go-server/pkg/service"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	_ "dubbo.apache.org/dubbo-go/v3/registry/xds"
)

// export DUBBO_GO_CONFIG_PATH=$PATH_TO_APP/conf/dubbogo.yaml
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
