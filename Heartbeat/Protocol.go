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
	//按照option所指设置发送request，并返回发回的信息和错误
	Send(request Request, option RequestOption, responseChan chan ResponseChanElement)
}

type ResponseOption interface {
}

//心跳数据响应协议
type ResponseProtocol interface {
	//接收到信息时将接收到的信息和错误入requestChan，并从responseChan中取出信息发回
	//
	//每收发一轮就调用一次此函数
	Recv(responseChan chan Response, option ResponseOption, requestChan chan RequestChanElement)
}
