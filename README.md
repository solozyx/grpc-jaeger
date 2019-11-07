# grpc-jaeger

grpc-jaeger is a kind of interceptor to gRPC implemented by Go, which is based on opentracing and uber/jaeger. You can use it to build a distributed gRPC-tracing system.

# Dependencies

```
github.com/opentracing/opentracing-go
github.com/uber/jaeger-client-go
```

# Testing

## get the result from the UI
Open the url `http://192.168.174.135:16686/`, and search the specified service. Then you should see the result like this as below：
![jaegerui](./jaegerui.png)