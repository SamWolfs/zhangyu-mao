ARG ALPINE_VERSION="latest"
ARG GO_VERSION="1.24.1"
FROM golang:$GO_VERSION AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY config config
COPY cmd cmd
COPY apps apps
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/zhangyumao.go

FROM alpine:$ALPINE_VERSION
WORKDIR /root/
COPY --from=builder /app/zhangyumao .
EXPOSE 3000
CMD ["./zhangyumao"]
