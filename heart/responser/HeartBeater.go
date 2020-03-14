package responser

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
)

//HeartBeater is the controller of the Heart.
//It is used in Heart to generate response.
type HeartBeater interface {
	//Beat was designed to generate a response (`protocol.Response` and `protocol.ResponseSendOption` in `protocol.TobeSendResponse`) from a received request (`Protocol.Request`).
	//In the implementaion of this method, you should do the following:
	//
	//1. Cast the `request` into your implementaion type
	//
	//2. Generate a `protocol.Response` you want to send back to the request and the send option `protocol.ResponseSendOption`
	//
	//3. Enclose the `protocol.Response` and `protocol.ResponseSendOption` into a `protocol.TobeSendResponse`
	//
	//4. Return the `protocol.TobeSendResponse`
	Beat(ctx context.Context, request protocol.Request) protocol.TobeSendResponse
}
