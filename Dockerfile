FROM registry.cn-hangzhou.aliyuncs.com/openfaas/golang:1.13-alpine3.11  AS builder
ENV GO111MODULE=on
ENV CGO_ENABLED 0
ENV GOOS=linux
ENV GOPROXY="https://goproxy.cn"

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -tags=jsoniter -ldflags="-s -w" -installsuffix cgo -o app main.go

FROM registry.cn-hangzhou.aliyuncs.com/epet/alpine:latest

COPY --from=builder /go/release/config /root/.kube/

COPY --from=builder /go/release/app /

CMD ["/app"]
