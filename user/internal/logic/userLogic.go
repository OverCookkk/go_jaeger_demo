package logic

import (
	"context"
	"go_jaeger_demo/pay/pay"

	"go_jaeger_demo/user/internal/svc"
	"go_jaeger_demo/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.Request) (resp *types.Response, err error) {
	l.Logger.Debugf("User name: %s", req.UserName)
	payResp, err := l.svcCtx.PayRpcClient.Greet(l.ctx, &pay.PayReq{Name: req.UserName})
	if err != nil {
		return nil, err
	}
	code := 0
	if payResp.Greet == "success" {
		code = 1
	}
	return &types.Response{Code: code}, nil
}
