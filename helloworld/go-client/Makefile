IMAGE = ccr.ccs.tencentyun.com/wumengye/dubbo-go-client
TAG = 1.0.1
HELM_INSTALL_NAME = dubbo-go-client

build-amd64-app:
	GOOS=linux GOARCH=amd64 go build -o build/app ./cmd

build: proto-gen tidy build-amd64-app
	cp ./conf/dubbogo.yaml ./build/dubbogo.yaml
	cp ./manifest/config/config.toml ./build/gf.toml
	docker build ./build -t ${IMAGE}:${TAG}
	docker push ${IMAGE}:${TAG}
	make clean

buildx-publish: proto-gen tidy build-amd64-app
	cp ./conf/dubbogo.yaml ./build/dubbogo.yaml
	cp ./manifest/config/config.toml ./build/gf.toml
	docker buildx build \
    	 --platform linux/amd64 \
    	 -t ${IMAGE}:${TAG} \
    	 ./build --push
	make clean

remove:
	helm uninstall ${HELM_INSTALL_NAME}

deploy:
	helm install ${HELM_INSTALL_NAME} ./chart/app

upgrade:
	helm upgrade ${HELM_INSTALL_NAME} ./chart/app

deploy-nacos-env:
	helm install nacos ./chart/nacos_env

remove-nacos-env:
	helm uninstall nacos

proto-gen:
	buf generate ../proto/helloworld.proto
	#protoc --go_out=./api --go-triple_out=./api ./api/api.proto

clean:
	rm ./build/dubbogo.yaml
	rm ./build/gf.toml
	rm ./build/app

tidy:
	go mod tidy

test:
	go test ./...

run: proto-gen tidy
	DUBBO_GO_CONFIG_PATH=conf/dubbogo-local.yaml go run cmd/*