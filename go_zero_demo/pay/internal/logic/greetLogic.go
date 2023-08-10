package logic

import (
	"context"

	"go_jaeger_demo/pay/internal/svc"
	"go_jaeger_demo/pay/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GreetLogic) Greet(in *pay.PayReq) (*pay.PayResp, error) {
	l.Logger.Debugf("PayReq name: %s", in.Name)
	if in.Name == "" {
		return &pay.PayResp{Greet: "failed"}, nil
	}
	return &pay.PayResp{Greet: "success"}, nil
}
