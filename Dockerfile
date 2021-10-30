FROM golang:1.17-alpine AS builder
WORKDIR /go/src/github.com/johnbellone/time-service
COPY . .
RUN apk add -U --no-cache ca-certificates git
RUN go mod tidy
RUN go build -o /go/bin/time-service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/time-service /bin/time-service
WORKDIR /root
EXPOSE 9090
ENTRYPOINT ["/bin/time-service"]
