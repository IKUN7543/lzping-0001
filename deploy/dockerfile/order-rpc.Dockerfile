FROM golang:1.22-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download 2>/dev/null || true

COPY . .

RUN go build -ldflags="-s -w" -o /app/order-rpc ./service/order/rpc/order.go

FROM alpine:3.19

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/order-rpc /app/order-rpc
COPY service/order/rpc/etc/*.yaml /app/etc/

EXPOSE 8031

CMD ["./order-rpc", "-f", "etc/order.yaml"]
