# time service

![GitHub Workflow](https://img.shields.io/github/workflow/status/johnbellone/time-service/go-workflow?style=for-the-badge)
[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](LICENSE)

An example gRPC service written in Go.

- Time Service for retrieving current and local system time.
- Multiple server versions running out of a single binary.
- Production level logging and tracing support.

## Usage

``` shell
➜  time-service-go git:(main) ✗ make
➜  time-service-go git:(main) ✗ bin/time-service &
[1] 94242
➜  time-service-go git:(main) ✗ {"level":"info","ts":1681314730.0254211,"caller":"time-service-go/main.go:143","msg":"build info","Executable":"/Users/JBellone/src/github.com/johnbellone/time-service-go/bin/time-service","Version":"0.1.0","GitAbbrv":"c51d3888-dirty","GitCommit":"c51d38883d03bb3eef2aab55efbdb116ac4c7756","BuildTime":"2023-04-12T11:52:01-0400"}
{"level":"info","ts":1681314730.0255392,"caller":"time-service-go/main.go:151","msg":"starting server","GrpcPort":50100}

➜  time-service-go git:(main) ✗ grpcurl -plaintext localhost:50100 list
grpc.health.v1.Health
grpc.reflection.v1alpha.ServerReflection
time.v1.TimeService
time.v2.TimeService
{"level":"info","ts":1681314737.4723892,"caller":"zap/options.go:212","msg":"finished streaming call with code OK","grpc.start_time":"2023-04-12T11:52:17-04:00","system":"grpc","span.kind":"server","grpc.service":"grpc.reflection.v1alpha.ServerReflection","grpc.method":"ServerReflectionInfo","peer.address":"[::1]:58582","grpc.code":"OK","grpc.time_ms":1.225}
➜  time-service-go git:(main) ✗ grpcurl -plaintext localhost:50100 time.v1.TimeService/GetCurrentTime
{"level":"info","ts":1681314750.922079,"caller":"zap/options.go:212","msg":"finished unary call with code OK","grpc.start_time":"2023-04-12T11:52:30-04:00","system":"grpc","span.kind":"server","grpc.service":"time.v1.TimeService","grpc.method":"GetCurrentTime","peer.address":"[::1]:58586","grpc.code":"OK","grpc.time_ms":0.009}
{
  "timestamp": "2023-04-12T15:52:30.922045Z"
}
{"level":"info","ts":1681314750.922843,"caller":"zap/options.go:212","msg":"finished streaming call with code OK","grpc.start_time":"2023-04-12T11:52:30-04:00","system":"grpc","span.kind":"server","grpc.service":"grpc.reflection.v1alpha.ServerReflection","grpc.method":"ServerReflectionInfo","peer.address":"[::1]:58586","grpc.code":"OK","grpc.time_ms":3.608}
➜  time-service-go git:(main) ✗ grpcurl -plaintext localhost:50100 time.v2.TimeService/GetCurrentTime
{"level":"info","ts":1681314758.660966,"caller":"zap/options.go:212","msg":"finished unary call with code OK","grpc.start_time":"2023-04-12T11:52:38-04:00","system":"grpc","span.kind":"server","grpc.service":"time.v2.TimeService","grpc.method":"GetCurrentTime","peer.address":"[::1]:58591","grpc.code":"OK","grpc.time_ms":0.011}
{
  "timestamp": "2023-04-12T15:52:38.660937Z"
}
{"level":"info","ts":1681314758.661383,"caller":"zap/options.go:212","msg":"finished streaming call with code OK","grpc.start_time":"2023-04-12T11:52:38-04:00","system":"grpc","span.kind":"server","grpc.service":"grpc.reflection.v1alpha.ServerReflection","grpc.method":"ServerReflectionInfo","peer.address":"[::1]:58591","grpc.code":"OK","grpc.time_ms":2.145}
➜  time-service-go git:(main) ✗ grpcurl -plaintext localhost:50100 time.v2.TimeService/GetLocalTime
{"level":"info","ts":1681314763.1889658,"caller":"zap/options.go:212","msg":"finished unary call with code OK","grpc.start_time":"2023-04-12T11:52:43-04:00","system":"grpc","span.kind":"server","grpc.service":"time.v2.TimeService","grpc.method":"GetLocalTime","peer.address":"[::1]:58594","grpc.code":"OK","grpc.time_ms":0.011}
{
  "timestamp": "2023-04-12T15:52:43.188940Z"
}
{"level":"info","ts":1681314763.189378,"caller":"zap/options.go:212","msg":"finished streaming call with code OK","grpc.start_time":"2023-04-12T11:52:43-04:00","system":"grpc","span.kind":"server","grpc.service":"grpc.reflection.v1alpha.ServerReflection","grpc.method":"ServerReflectionInfo","peer.address":"[::1]:58594","grpc.code":"OK","grpc.time_ms":2.93}
```

## License

`time-service` is released under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
