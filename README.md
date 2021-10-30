# time service

![GitHub Workflow](https://img.shields.io/github/workflow/status/johnbellone/time-service/go-workflow?style=for-the-badge)
[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](LICENSE)

An example gRPC service written in Go.

## Usage

``` shell
~ % grpcurl -insecure localhost:9090 list
grpc.health.v1.Health
grpc.reflection.v1alpha.ServerReflection
time.v1.Time
~ % grpcurl -insecure localhost:9090 describe time.v1.Time
time.v1.Time is a service:
service Time {
  rpc GetCurrentTime ( .time.v1.TimeRequest ) returns ( .time.v1.TimeResponse ) {
    option (.google.api.http) = { get:"/v1/time"  };
  }
}
~ % grpcurl -insecure localhost:9090 time.v1.Time/GetCurrentTime
{
  "currentTime": "2021-10-30T13:51:45.033413Z"
}
```

## License

`time-service` is released under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
