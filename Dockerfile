FROM golang:1.24.0-alpine3.21 AS builder

ARG SHORT_COMMIT=development
ARG VERSION=development
ARG BUILD_TIME=unknown

WORKDIR /app
COPY . .

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build \
       -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.shortCommit=${SHORT_COMMIT}" \
       -a -installsuffix cgo \
       -o order-api-service .

FROM alpine:3.21

ENV TZ=Asia/Jakarta

RUN apk add --no-cache ca-certificates tzdata \
    && ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime \
    && echo "$TZ" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/order-api-service .
RUN mkdir -p configs
COPY --from=builder /app/configs/application.sample.yaml ./configs/application.yaml

RUN chmod +x order-api-service

USER 1001

ENTRYPOINT ["./order-api-service", "server"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3002/ping || exit 1