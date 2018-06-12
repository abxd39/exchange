1.gatewayw服务器负责网关
（1）conf组件包括目录配置文件
（2）http组件对外部暴漏的基于gin的web接口，controller层进行对应逻辑处理
（3）rpc组件进行内部rpc调用通过web接口触发（rpc客户端）
（4）gateway是提供对外访问程序入口

2.user_service是单独的一个处理用户逻辑的微服务
（1）conf组件包括目录配置文件
（2）dao组件是对mysql数据库和redis的封装，使用对应接口进行数据层操作
（3）rpc组件是一个提供对用户数据处理的对外暴漏的接口（rpc服务器）
（4）model提供内部数据结构统一定义管理
（5）user.go是用户功能入口

3.proto文件夹是定义内部通信协议目录
(1)现提供一个网关到用户服务消息定义文件gateway2user.proto按照标准protobuf服务格式定义。
提供一个hello接口实例
service Gateway2User {
	rpc Hello(HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
	string name = 1;
}

message HelloResponse {
	string greeting = 2;
}

####
使用micro框架开发接口流程
1.在proto文件夹生成所需要的内部调用接口(执行脚本自动生成sh run.sh)
2.在user_service/rpc/handler中新增内部rpc调用的hello处理函数func (s *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error
3.在gateway/rpc/client/中新增所需要调用的服务类型（实现实例所调用的user 的CallGreet方法）
4.在http/controller层新增对应的处理方法HelloController(c *gin.Context)
5.先在gateway服务器http/http.go文件增加路由r.GET("/test/:name", controller.HelloController)



