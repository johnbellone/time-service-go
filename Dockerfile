FROM golang:1.20-bullseye AS builder
WORKDIR /go/src/github.com/johnbellone/time-service-go
COPY . .
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
    ; \ 
    rm -rf /var/lib/apt/lists/*
RUN go mod tidy
RUN go build -o /go/bin/time-service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/time-service /bin/time-service
WORKDIR /root
ENTRYPOINT ["/bin/time-service"]
