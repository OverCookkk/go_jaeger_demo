package logic

import (
	"context"
	"go_jaeger_demo/order/internal/svc"
	"go_jaeger_demo/order/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderLogic) Order(req *types.Request) (resp *types.Response, err error) {
	l.Logger.Debugf("OrderReq name: %s", req.OrderName)
	if req.OrderName == "" {
		return &types.Response{Msg: "failed"}, nil
	}
	return &types.Response{Msg: "success"}, nil
}
