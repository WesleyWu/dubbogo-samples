FROM alpine:3.15

WORKDIR /app

ENV DUBBO_GO_CONFIG_PATH=dubbogo.yaml
ENV GF_GCFG_FILE=gf.toml

COPY ./app /app/app
COPY ./dubbogo.yaml /app/dubbogo.yaml
COPY ./gf.toml /app/gf.toml

CMD /app/app
