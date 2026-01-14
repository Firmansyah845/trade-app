FROM golang:1.24.0-alpine3.21 AS builder

ARG SHORT_COMMMIT=development

WORKDIR /app
COPY . .

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.shortCommit=${SHORT_COMMMIT}" -a -installsuffix cgo -o order-api-service .

FROM alpine:3.20

ENV TZ=Asia/Jakarta

RUN apk update && apk --no-cache add ca-certificates bash jq curl tzdata \
    && rm -rf /var/cache/apk/* \
    && ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime \
    && echo "$TZ" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/ads-api-service .
RUN mkdir -p configs
COPY --from=builder /app/configs/application.sample.yaml ./configs/application.yaml

RUN chmod +x ads-api-service
USER 1001

ENTRYPOINT ["./warehouse-api-service ", "server"]

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=3 CMD curl --fail http://localhost:3002/ping
