package Heart

type RequesterHeartbeat interface{}
type ResponserHeartbeat interface{}

type RequesterHeartbeatProtocol interface {
	//启下：对接下层协议
	//
	//发送一个Heartbeat数据请求，并返回响应和错误
	RequestHeartbeat(beat RequesterHeartbeat) (ResponserHeartbeat, error)

	//承上：对接上层消息策略
	//
	//输入一个Heartbeat数据响应和下一个Heartbeat处理函数，处理响应并生成下一个Heartbeat数据请求
	Beat(request RequesterHeartbeat, response ResponserHeartbeat, beat func(RequesterHeartbeat))
}
