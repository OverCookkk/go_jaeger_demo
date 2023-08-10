package svc

import (
	"context"
	"github.com/zeromicro/go-zero/zrpc"
	"go_jaeger_demo/pay/payclient"
	"go_jaeger_demo/user/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServiceContext struct {
	Config       config.Config
	PayRpcClient payclient.Pay
	// OrderRpcClient orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		PayRpcClient: payclient.NewPay(zrpc.MustNewClient(c.PayRpcConf, zrpc.WithUnaryClientInterceptor(TokenuidIntereptor))),
		// OrderRpcClient: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpcConf, zrpc.WithUnaryClientInterceptor(TokenuidIntereptor))),
	}
}

// TokenuidIntereptor 函数是一个 gRPC 客户端拦截器
func TokenuidIntereptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	md := metadata.New(map[string]string{"username": "zhangsan"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	err = invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}
	return nil
}
