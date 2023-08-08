# go_jaeger_demo

基于go-zero框架写的jaeger链路追踪demo，包含两个服务，user和order是一个提供http接口的服务，pay是一个提供rpc接口的服务，具体调用方式如图所示：

![go_jaeger_demo_2](https://raw.githubusercontent.com/OverCookkk/PicBed/master/blogImg/go_jaeger_demo_2.jpg)

## 部署与运行

### jaeger服务的部署

下载地址：https://www.jaegertracing.io/download/

解压文件后运行服务：

```
./jaeger-all-in-one.exe
```

启动后，访问地址：http://localhost:16686/，界面如下：

![go_jaeger_demo_1](https://raw.githubusercontent.com/OverCookkk/PicBed/master/blogImg/go_jaeger_demo_1.png)





### 启动服务

go-zero框架本身就 已经支持jaeger，只需要配置好就可以

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







## 后续

- [ ] 外部非go-zero框架的服务调用go-zero框架的api和rpc服务貌似不能直接用http头传递trace_id，需要加入propagation注入，有待研究。
