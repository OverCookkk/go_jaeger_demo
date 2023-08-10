package main

import (
	"bufio"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"io"
	"net/http"
	"os"
	"time"
)

func CreateTracer(servieName string) (opentracing.Tracer, io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: servieName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 按实际情况替换你的 ip
			CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
			// LocalAgentHostPort: "192.168.129.128:31300",
		},
	}

	jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
	)
	return tracer, closer, err
}

func main() {
	tracer, closer, _ := CreateTracer("UserinfoService")
	// 创建第一个 span A
	parentSpan := tracer.StartSpan(
		"clientmain",
		opentracing.StartTime(time.Now()),
	)

	// 调用其它服务
	GetUserInfo(tracer, parentSpan)
	// 结束 A
	parentSpan.Finish()
	// 结束当前 tracer
	closer.Close()

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadByte()
}

// 请求远程服务，获得用户信息
func GetUserInfo(tracer opentracing.Tracer, parentSpan opentracing.Span) {
	// 继承上下文关系，创建子 span
	childSpan := tracer.StartSpan(
		"GetUserInfoRequest",
		opentracing.ChildOf(parentSpan.Context()),
	)

	url := "http://127.0.0.1:8081/Get?username=hzh"
	req, _ := http.NewRequest("GET", url, nil)
	// 设置 tag
	ext.SpanKindRPCClient.Set(childSpan)
	childSpan.SetTag("myself tag", "tag123")
	ext.SpanKind.Set(childSpan, "client")
	ext.HTTPUrl.Set(childSpan, url)
	ext.HTTPMethod.Set(childSpan, "GET")
	childSpan.LogFields(log.String("event", "a test event"))
	// 在http的header中注入追踪信息，以便在服务端Extract()获取追踪信息
	tracer.Inject(childSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)

	_ = resp                       // 丢掉
	ext.Error.Set(childSpan, true) // 创建一个error信息
	defer childSpan.Finish()
}
