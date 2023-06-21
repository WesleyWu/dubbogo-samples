IMAGE = ccr.ccs.tencentyun.com/wumengye/dubbo-go-server
TAG = 1.0.6
SERVICE_NAME = dubbogo-helloworld
SERVICE_GRPC_PORT = 20000
HELM_INSTALL_NAME = ${SERVICE_NAME}-server-v1
HEALTH_CHECK_GLOBAL = grpc-gke-hc

gcp-firewall-rules:
	gcloud compute firewall-rules describe ${SERVICE_NAME}-allow-hc    \
        && echo "firewall-rules ${SERVICE_NAME}-allow-hc already exists" || \
        gcloud compute firewall-rules create ${SERVICE_NAME}-allow-hc \
          --network default --action allow --direction INGRESS \
          --source-ranges 35.191.0.0/16,130.211.0.0/22 \
          --target-tags allow-health-checks \
          --rules tcp:${SERVICE_GRPC_PORT}

gcp-backend-services:
	gcloud compute backend-services describe ${SERVICE_NAME}-server \
      --global \
        && echo "backend-services ${SERVICE_NAME}-server already exists" || \
    gcloud compute backend-services create ${SERVICE_NAME}-server \
      --global \
      --load-balancing-scheme=INTERNAL_SELF_MANAGED \
      --protocol=GRPC \
      --health-checks ${HEALTH_CHECK_GLOBAL}

HAS_BACKEND = $(shell gcloud compute backend-services describe ${SERVICE_NAME}-server --global | grep ${SERVICE_NAME}-service-neg | wc -c | xargs)
gcp-backend-services-add-backend:gcp-backend-services
	@if [ "$(HAS_BACKEND)" -eq "0" ]; then \
		gcloud compute backend-services add-backend ${SERVICE_NAME}-server \
			--global \
			--network-endpoint-group ${SERVICE_NAME}-service-neg \
			--network-endpoint-group-zone us-central1-a \
			--balancing-mode RATE \
			--max-rate-per-endpoint 5; \
    fi

gcp-url-maps:
	gcloud compute url-maps describe ${SERVICE_NAME}-url-map \
        && echo "url-maps ${SERVICE_NAME}-url-map already exists" || \
	gcloud compute url-maps create ${SERVICE_NAME}-url-map \
      --default-service ${SERVICE_NAME}-server

HAS_PATH_MATCHER = $(shell gcloud compute url-maps describe ${SERVICE_NAME}-url-map | grep ${SERVICE_NAME}-path-matcher | wc -c | xargs)
gcp-url-maps-path-matcher:gcp-url-maps
	@if [ "$(HAS_PATH_MATCHER)" -eq "0" ]; then \
		gcloud compute url-maps add-path-matcher ${SERVICE_NAME}-url-map \
			  --default-service ${SERVICE_NAME}-server \
			  --path-matcher-name ${SERVICE_NAME}-path-matcher \
			  --new-hosts ${SERVICE_NAME}:${SERVICE_GRPC_PORT}; \
	fi

gcp-target-grpc-proxies:gcp-url-maps
	gcloud compute target-grpc-proxies describe ${SERVICE_NAME}-proxy \
        && echo "target-grpc-proxies ${SERVICE_NAME}-proxy already exists" || \
	gcloud compute target-grpc-proxies create ${SERVICE_NAME}-proxy \
      --url-map ${SERVICE_NAME}-url-map \
      --validate-for-proxyless

gcp-forwarding-rules:gcp-target-grpc-proxies
	gcloud compute forwarding-rules describe ${SERVICE_NAME}-forwarding-rule --global \
		&& echo "forwarding-rules ${SERVICE_NAME}-forwarding-rule already exists" || \
	gcloud compute forwarding-rules create ${SERVICE_NAME}-forwarding-rule \
          --global \
          --load-balancing-scheme=INTERNAL_SELF_MANAGED \
          --address=0.0.0.0 \
          --target-grpc-proxy=${SERVICE_NAME}-proxy \
          --ports ${SERVICE_GRPC_PORT} \
          --network default

gcp:gcp-firewall-rules gcp-backend-services-add-backend gcp-url-maps-path-matcher gcp-forwarding-rules
	@echo "GCP traffic director configured for service ${SERVICE_NAME}"

gcp-clear:
	gcloud compute forwarding-rules delete ${SERVICE_NAME}-forwarding-rule --global -q
	echo $(gcloud compute target-grpc-proxies delete ${SERVICE_NAME}-proxy -q)
	echo $(gcloud compute url-maps delete ${SERVICE_NAME}-url-map -q)
	echo $(gcloud compute backend-services delete ${SERVICE_NAME}-server --global -q)
	echo $(gcloud compute firewall-rules delete ${SERVICE_NAME}-allow-hc -q)
	echo $(gcloud compute network-endpoint-groups delete ${SERVICE_NAME}-service-neg --zone us-central1-a -q)
	echo "done"

build-amd64-app:
	GOOS=linux GOARCH=amd64 go build -o build/app ./cmd

build:proto-gen tidy build-amd64-app
	cp ./conf/dubbogo.yaml ./build/dubbogo.yaml
	docker build ./build -t ${IMAGE}:${TAG}
	docker push ${IMAGE}:${TAG}
	make clean

buildx-publish:proto-gen tidy build-amd64-app
	cp ./conf/dubbogo.yaml ./build/dubbogo.yaml
	docker buildx build \
         --platform linux/amd64 \
         -t ${IMAGE}:${TAG} \
         ./build --push
	make clean

remove-gke:
	helm uninstall ${HELM_INSTALL_NAME}

remove:
	helm uninstall ${HELM_INSTALL_NAME}

deploy-gke:
	helm install ${HELM_INSTALL_NAME} ./chart/gke

deploy:
	helm install ${HELM_INSTALL_NAME} ./chart/app

upgrade-gke:
	helm upgrade ${HELM_INSTALL_NAME} ./chart/gke

upgrade:
	helm upgrade ${HELM_INSTALL_NAME} ./chart/app

deploy-nacos-env:
	helm install nacos ./chart/nacos_env

remove-nacos-env:
	helm uninstall nacos

proto-gen:
	buf generate ../proto/helloworld.proto

clean:
	rm ./build/dubbogo.yaml
	rm ./build/app

tidy:
	go mod tidy

test:
	go test ./...

run:proto-gen tidy
	DUBBO_GO_CONFIG_PATH=conf/dubbogo-local.yaml go run cmd/*