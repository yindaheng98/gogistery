package Heartbeat

//心跳数据请求基础类
type Request interface {
}

//心跳数据响应基础类
type Response interface {
}

//此类用于在Response通道中传递Response和错误信息
type ResponseChanElement struct {
	response Response
	error    error
}

func (e ResponseChanElement) GetResponse() Response {
	return e.response
}

func (e ResponseChanElement) GetError() error {
	return e.error
}

//此类用于在Request通道中传递Request和错误信息
type RequestChanElement struct {
	request Request
	error   error
}

type RequestOption interface {
}

//心跳数据发送协议
type RequestProtocol interface {
	//按照option所指设置从只读channel responseChan中取出信息发出，并将发回的信息和错误放入只写channel responseChan
	Request(request <-chan Request, option RequestOption, responseChan chan<- ResponseChanElement)
}

type ResponseOption interface {
}

//心跳数据响应协议
type ResponseProtocol interface {
	//接收到信息时将接收到的信息和错误放入只写channel requestChan，并从只读channel responseChan中取出信息发回
	Response(requestChan chan<- RequestChanElement, option ResponseOption, responseChan <-chan Response)
}
