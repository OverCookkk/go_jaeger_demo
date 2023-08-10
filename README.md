# go_jaeger_demo

该项目包含两种类型的go使用jaeger的demo，分别存在**go_zero_demo**和**other_demo**两个目录。

## 部署与运行

### jaeger服务的部署

下载地址：https://www.jaegertracing.io/download/

解压文件后运行服务：

```
./jaeger-all-in-one.exe
```

启动后，访问地址：http://localhost:16686/，界面如下：

![go_jaeger_demo_1](https://raw.githubusercontent.com/OverCookkk/PicBed/master/blogImg/go_jaeger_demo_1.png)





### 基于go-zero框架的demo

基于go-zero框架写的jaeger链路追踪demo，**go_zero_demo**目录包含三个服务，user和order是一个提供http接口的服务，pay是一个提供rpc接口的服务，具体调用方式如图所示：

<img src="https://raw.githubusercontent.com/OverCookkk/PicBed/master/blogImg/go_jaeger_demo_2.jpg" alt="go_jaeger_demo_2" style="zoom:50%;" />



go-zero框架本身就已经支持jaeger，只需要配置好就可以；

**注意：go-zero框架的api和rpc服务调用下游服务的使用都需要传递上下文context，这样下游服务和上游服务才会是同一个trace_id**

1. 在yaml配置上加入jaeger配置

   ```yaml
   Telemetry:
     Name: order.api
     Endpoint: http://127.0.0.1:14268/api/traces
     Sampler: 1.0
     Batcher: jaeger
   ```

2. 不管是HTTP还是RPC通信，在调用下游服务的时候，把上下文context传递出去，然后启动user、pay、order三个服务。

3. 请求user服务后，再在jaeger UI上查看调用信息

   ![go_jaeger_demo_3](https://raw.githubusercontent.com/OverCookkk/PicBed/master/blogImg/go_jaeger_demo_3.png)





### 基于其他框架的demo

**other_demo**目录包含三个服务，user和order是一个提供http接口的服务，pay是一个提供rpc接口的服务，主要使用`"github.com/opentracing/opentracing-go"`和`"github.com/uber/jaeger-client-go"`库。

1. 对于http调用，客户端生成的追踪信息注入http的header中，服务端再从http的header中获取出来，从而把链路关联起来。

2. 对于grpc使用3方库来操作，调研测试了grpc-jaeger，grpc-jaeger 是 Go 实现的一种 gRPC 拦截器，它基于 opentracing 和 uber/jaeger。您可以使用它来构建分布式 gRPC 跟踪系统。它的内部不需要特殊处理就可以把跨进程的调用串联到一起。





## 后续

- [ ] 外部非go-zero框架的服务调用go-zero框架的api和rpc服务貌似不能直接用http头传递trace_id，需要加入propagation注入，有待研究。
