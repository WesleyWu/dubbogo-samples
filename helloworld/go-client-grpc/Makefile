IMAGE = ccr.ccs.tencentyun.com/wumengye/dubbo-go-client-grpc
TAG = 1.0.13
HELM_INSTALL_NAME = dubbo-go-client-grpc

build-amd64-app:
	GOOS=linux GOARCH=amd64 go build -o build/app main.go

build: proto-gen tidy build-amd64-app
	cp ./manifest/config/config.toml ./build/gf.toml
	cp ./xds_bootstrap.json ./build
	docker build ./build -t ${IMAGE}:${TAG}
	docker push ${IMAGE}:${TAG}
	make clean

remove:
	helm uninstall ${HELM_INSTALL_NAME}

deploy:
	helm install ${HELM_INSTALL_NAME} ./chart/gke

upgrade:
	helm upgrade ${HELM_INSTALL_NAME} ./chart/gke

proto-gen:
	buf generate ../proto/helloworld.proto
	#protoc --go_out=./api --go-triple_out=./api ./api/api.proto

clean:
	rm ./build/gf.toml
	rm ./build/app
	rm ./build/xds_bootstrap.json

tidy:
	go mod tidy

test:
	go test ./...

run: proto-gen tidy
	go run cmd/*