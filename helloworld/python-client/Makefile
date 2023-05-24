IMAGE = ccr.ccs.tencentyun.com/wumengye/dubbo-go-python-client
TAG = 1.0.11
HELM_INSTALL_NAME = dubbo-go-python-client

install:
	pip3 install -r requirements.txt

build: proto-gen
	cp ./*.py ./build
	cp ./requirements.txt ./build
	cp ./xds_bootstrap.json ./build
	cp -r ./proto ./build
	docker build ./build -t ${IMAGE}:${TAG}
	docker push ${IMAGE}:${TAG}
	make clean

remove:
	helm uninstall ${HELM_INSTALL_NAME}

deploy:
	helm install ${HELM_INSTALL_NAME} ./chart/app

upgrade:
	helm upgrade ${HELM_INSTALL_NAME} ./chart/app

proto-gen:
	buf generate ../proto/helloworld.proto

run: proto-gen
#	SERVICE_ENDPOINT=127.0.0.1:20000 uvicorn main:app --host 0.0.0.0 --port 5000
	SERVICE_ENDPOINT=192.168.31.220:20000 uvicorn main:app --host 0.0.0.0 --port 5000

clean:
	rm ./build/main.py
	rm ./build/requirements.txt
	rm ./build/xds_bootstrap.json
	rm -r ./build/proto
