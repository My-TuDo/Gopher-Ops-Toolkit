FROM golang:1.26-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/gopt .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/gopt /gopt

ENTRYPOINT ["/gopt"]
CMD ["--help"]