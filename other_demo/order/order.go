package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"io"
	"net/http"
)

func main() {
	r := gin.Default()
	// 插入中间件处理
	r.Use(UseOpenTracing())
	r.GET("/Get", GetUserInfo)
	r.Run("127.0.0.1:8081")
}

// 从上下文中解析并创建一个新的 trace，获得传播的 上下文(SpanContext)
func CreateTracer(serviceName string, header http.Header) (opentracing.Tracer, opentracing.SpanContext, io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 上报给agent
			// LocalAgentHostPort: "192.168.128.128:31300",
			// 上报给collector
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
		},
	}

	jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
	)
	// 继承别的进程传递过来的上下文，对应客户端的tracer.Inject()操作
	spanContext, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
	return tracer, spanContext, closer, err
}

func UseOpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 opentracing.GlobalTracer() 获取全局 Tracer
		tracer, spanContext, closer, _ := CreateTracer("userInfoWebService", c.Request.Header)
		defer closer.Close()
		// 生成依赖关系，并新建一个 span
		// 这里很重要，因为生成了  References []SpanReference 依赖关系
		startSpan := tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(spanContext))
		defer startSpan.Finish()

		// 记录 tag
		// 记录请求 Url
		ext.HTTPUrl.Set(startSpan, c.Request.URL.Path)
		// Http Method
		ext.HTTPMethod.Set(startSpan, c.Request.Method)
		// 记录组件名称
		ext.Component.Set(startSpan, "Gin-Http")

		// 在 header 中加上当前进程的上下文信息
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), startSpan))
		// 传递给下一个中间件
		c.Next()
		// 继续设置 tag
		ext.HTTPStatusCode.Set(startSpan, uint16(c.Writer.Status()))
	}
}

func GetUserInfo(ctx *gin.Context) {
	userName := ctx.Param("username")
	fmt.Println("收到请求，用户名称为:", userName)
	ctx.String(http.StatusOK, "hello world")
}
