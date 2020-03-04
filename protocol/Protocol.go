package protocol

import (
	"context"
	"fmt"
)

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

//发送一个请求所需的信息
type TobeSendRequest struct {
	Request Request
	Option  RequestSendOption
}

func (r TobeSendRequest) String() string {
	return fmt.Sprintf("TobeSendRequest{Request:%s,Option:%s}", r.Request.String(), r.Option.String())
}

//发送一个响应所需的信息
type TobeSendResponse struct {
	Response Response
	Option   ResponseSendOption
}

func (r TobeSendResponse) String() string {
	return fmt.Sprintf("TobeSendResponse{Response:%s,Option:%s}", r.Response.String(), r.Option.String())
}

//心跳数据发送协议
type RequestProtocol interface {
	//从只读channel responseChan中取出信息发出，并将发回的信息和错误放入只写channel responseChan
	Request(ctx context.Context, requestChan <-chan TobeSendRequest, responseChan chan<- ReceivedResponse)
}

//心跳数据响应协议
type ResponseProtocol interface {
	//接收到信息时将接收到的信息和错误放入只写channel requestChan，并从只读channel responseChan中取出信息发回
	Response(ctx context.Context, requestChan chan<- ReceivedRequest, responseChan <-chan TobeSendResponse)
}
