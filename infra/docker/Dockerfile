FROM golang:1.23-bullseye AS build

WORKDIR /app

COPY . .

RUN go test -v ./... && go build -ldflags="-w -s" -o /app/etcd-ui internal/cmd/main.go

FROM alpine:3.18.2

WORKDIR /app
COPY --from=build --link /app/etcd-ui /app/etcd-ui

ENV LOG_DIR_PATH=/home/etcdui/log/

RUN adduser -D etcdui && \
    mkdir /home/etcdui/log && \
    chown etcdui:etcdui -R /home/etcdui/ /app && \
    apk add gcompat --no-cache

USER etcdui
EXPOSE 8080/tcp

CMD ["/app/etcd-ui"]
