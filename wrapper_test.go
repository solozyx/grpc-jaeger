package wrapper

import (
	"fmt"
	"net"
	"os"
	"time"

	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/solozyx/grpc-jaeger/proto"
)

var (
	JaegerHost = "192.168.174.135:6831"

	GrpcServer = "127.0.0.1:8001"
)

func Test_Tracing(t *testing.T) {
	ln, err := net.Listen("tcp", GrpcServer)
	if err != nil {
		os.Exit(-1)
	}

	var servOpts []grpc.ServerOption
	tracer, _, err := NewJaegerTracer("grpc-jaeger-server", JaegerHost)
	if err != nil {
		fmt.Printf("new tracer err: %+v\n", err)
		os.Exit(-1)
	}
	if tracer != nil {
		servOpts = append(servOpts, ServerOption(tracer))
	}
	srv := grpc.NewServer(servOpts...)
	pb.RegisterGreeterSrv(srv)

	go func() {
		time.Sleep(time.Second)

		dialOpts := []grpc.DialOption{grpc.WithInsecure()}
		tracer, _, err := NewJaegerTracer("grpc-jaeger-client", JaegerHost)
		if err != nil {
			fmt.Printf("new tracer err: %+v\n", err)
			os.Exit(-1)
		}
		if tracer != nil {
			dialOpts = append(dialOpts, DialOption(tracer))
		}
		conn, err := grpc.Dial(GrpcServer, dialOpts...)
		if err != nil {
			fmt.Printf("grpc connect failed, err:%+v\n", err)
			os.Exit(-1)
		}
		defer conn.Close()

		client := pb.NewGreeterClient(conn)
		reqBody := pb.HelloRequest{
			Name:    "i am tester",
			Message: "just for test",
		}
		resp, err := client.SayHello(context.Background(), &reqBody)
		if err != nil {
			fmt.Printf("call sayhello failed, err:%+v\n", err)
			os.Exit(-1)
		} else {
			fmt.Printf("call sayhello suc, res:%+v\n", resp)
		}
	}()

	go srv.Serve(ln)

	time.Sleep(time.Second * 5)
}
