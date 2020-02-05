package Heartbeat

//心跳数据请求基础类
type Request interface {
}

//心跳数据响应基础类
type Response interface {
}

//此类用于存储request端收到的response和错误信息
type ReceivedResponse struct {
	Response Response
	Error    error
}

//此类用于存储response端收到的request和错误信息
type ReceivedRequest struct {
	Request Request
	Error   error
}

//自定义请求发送设置
type RequestSendOption interface {
}

//自定义响应发送设置
type ResponseSendOption interface {
}

//发送一个请求所需的信息
type TobeSendRequest struct {
	Request Request
	Option  RequestSendOption
}

//发送一个响应所需的信息
type TobeSendResponse struct {
	Response Response
	Option   ResponseSendOption
}

//心跳数据发送协议
type RequestBeatProtocol interface {
	//从只读channel responseChan中取出信息发出，并将发回的信息和错误放入只写channel responseChan
	Request(requestChan <-chan TobeSendRequest, responseChan chan<- ReceivedResponse)
}

//心跳数据响应协议
type ResponseBeatProtocol interface {
	//接收到信息时将接收到的信息和错误放入只写channel requestChan，并从只读channel responseChan中取出信息发回
	Response(requestChan chan<- ReceivedRequest, responseChan <-chan TobeSendResponse)
}
