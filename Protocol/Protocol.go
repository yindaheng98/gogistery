package Protocol

import "gogistery/Controller"

//此类用于在Requester中传递Responser回传的信息或发送失败传回的错误信息
type RequestChanElement struct {
	Request Controller.Request
	Error   error
}

type Requester interface {
	//将request发送到addr去，并将返回的信息和错误入responseChan
	//
	//每一次发送都调用一次此函数。如果超时，responseChan将直接关闭，请自行处理超时panic
	Send(request Controller.Request, addr string, responseChan chan RequestChanElement)
}

//此类用于在Responser中传递接收到的信息或接收失败传回的错误信息
type ResponseChanElement struct {
	Response Controller.Response
	Error    error
}

type Responser interface {
	//接收到信息时将接收到的信息和错误入requestChan，并从responseChan中取出信息发回
	//
	//每收发一轮就调用一次此函数
	Recv(requestChan chan RequestChanElement, responseChan chan Controller.Response)
}
