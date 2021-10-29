FROM golang:1.16-alpine AS builder
WORKDIR /go/src/git.jbellone.dev/xchg/go-service
COPY . .
RUN apk add -U --no-cache ca-certificates git
RUN go mod tidy
RUN go build -o /go/bin/go-service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/go-service /bin/go-service
WORKDIR /root
EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/go-service"]
