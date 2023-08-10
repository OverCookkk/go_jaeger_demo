package main

import (
    "context"
    "fmt"
    "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/log"
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
    "io"
)

func Init(service string) (opentracing.Tracer, io.Closer) {
    // trace 配置
    cfg := &config.Configuration{
        ServiceName: service,
        Sampler: &config.SamplerConfig{
            Type:  jaeger.SamplerTypeConst,
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans: true,
            // collector 信息根据自己ip配置
            CollectorEndpoint: "http://127.0.0.1:14268/api/traces",
        },
    }
    // 根据上面的配置新建一个tracer
    tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
    if err != nil {
        panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
    }
    return tracer, closer
}
func main() {
    tracer, closer := Init("hello-world")
    helloTo := "rookie in jaeger"
    defer closer.Close()

    opentracing.SetGlobalTracer(tracer)

    // 创建一个span并且设置tag
    span := tracer.StartSpan("say-hello")
    span.SetTag("hello-to", helloTo)

    ctx := opentracing.ContextWithSpan(context.Background(), span)

    helloStr := formatString(ctx, helloTo)
    printHello(ctx, helloStr)

    // helloStr := fmt.Sprintf("Hello, %s!", helloTo)
    // // LogFields和LogKV都可以设置log
    // span.LogFields(
    // 	log.String("event", "string-format"),
    // 	log.String("value", helloStr),
    // )
    //
    // println(helloStr)
    // span.LogKV("event", "println")

    span.Finish()

}

func formatString(ctx context.Context, helloTo string) string {
    span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
    defer span.Finish()

    helloStr := fmt.Sprintf("Hello, %s!", helloTo)
    span.LogFields(
        log.String("event", "string-format"),
        log.String("value", helloStr),
    )

    return helloStr
}

func printHello(ctx context.Context, helloStr string) {
    span, _ := opentracing.StartSpanFromContext(ctx, "printHello")
    defer span.Finish()

    println(helloStr)
    span.LogKV("event", "println")
}
