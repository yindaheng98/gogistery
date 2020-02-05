package Heart

type RequesterBeat interface{}
type ResponserBeat interface{}
type RequesterBeatSendOption interface{}
type ResponserBeatSendOption interface{}
type TobeSendRequesterBeat struct {
	RequesterBeat RequesterBeat
	SendOption    RequesterBeatSendOption
}
type TobeSendResponserBeat struct {
	ResponserBeat ResponserBeat
	SendOption    ResponserBeatSendOption
}

type RequesterHeartProtocol interface {
	//启下：对接下层协议
	//
	//发送一个Beat数据请求，并返回响应和错误
	Request(beat TobeSendRequesterBeat) (ResponserBeat, error)

	//承上：对接上层消息策略
	//
	//输入一个Beat数据响应和下一个Beat处理函数，处理响应并生成下一个Beat数据请求
	Beat(request TobeSendRequesterBeat, response ResponserBeat, beat func(TobeSendRequesterBeat))
}

type ResponserHeartProtocol interface {
	//启下：对接下层协议
	//
	//接收一个Beat数据请求，并从响应队列中取出响应发回
	Response() (RequesterBeat, error, func(TobeSendResponserBeat))

	//承上：对接上层消息策略，每一个成功到达的数据请求都必须有响应
	//
	//输入一个Beat数据请求，处理请求并生成Beat数据响应
	Beat(request RequesterBeat) TobeSendResponserBeat
}
