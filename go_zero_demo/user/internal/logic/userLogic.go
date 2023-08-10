package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
	"go_jaeger_demo/pay/pay"
	"go_jaeger_demo/user/internal/svc"
	"go_jaeger_demo/user/internal/types"
	"io/ioutil"
	"net/http"
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

	type hRequest struct {
		// Node   string `path:"node"`
		OrderName string `json:"order_name"`
		// Header    string `header:"X-Header-Trace"`
	}
	hreq := hRequest{
		// Header:    trace.TraceIDFromContext(l.ctx),
		OrderName: req.UserName,
	}
	// 把上下文（l.ctx）传递给下游服务，让下游服务也拥有相同的trace_id
	hresp, err := httpc.Do(l.ctx, http.MethodPost, "http://127.0.0.1:9003/order", hreq)
	if err != nil {
		l.Logger.Debug(err)
		return
	}
	bodyBytes, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		// 处理错误
		return
	}
	defer hresp.Body.Close()

	l.Logger.Debugf("bodyString: %s", string(bodyBytes))

	// orderResp, err := l.svcCtx.OrderRpcClient.CreateOrder(l.ctx, &order.OrderReq{Name: req.UserName})
	// if err != nil {
	//     return nil, err
	// }
	// if orderResp.Code == "order success" {
	//     code = 2
	// }
	return &types.Response{Code: code}, nil
}
