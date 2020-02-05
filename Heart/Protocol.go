package Heart

type RequesterHeartbeat interface{}
type ResponserHeartbeat interface{}

type RequesterHeartProtocol interface {
	//启下：对接下层协议
	//
	//发送一个Heartbeat数据请求，并返回响应和错误
	Request(beat RequesterHeartbeat) (ResponserHeartbeat, error)

	//承上：对接上层消息策略
	//
	//输入一个Heartbeat数据响应和下一个Heartbeat处理函数，处理响应并生成下一个Heartbeat数据请求
	Beat(request RequesterHeartbeat, response ResponserHeartbeat, beat func(RequesterHeartbeat))
}

type ResponserHeartProtocol interface {
	//启下：对接下层协议
	//
	//接收一个Heartbeat数据请求，并从响应队列中取出响应发回
	Response() (RequesterHeartbeat, error, func(ResponserHeartbeat))

	//承上：对接上层消息策略，每一个成功到达的数据请求都必须有响应
	//
	//输入一个Heartbeat数据请求，处理请求并生成Heartbeat数据响应
	Beat(request RequesterHeartbeat) ResponserHeartbeat
}
