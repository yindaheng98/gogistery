package requester

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//This interface is used in Heart to generate request.
//In `Heart.RunBeating(...)`, when sent a request and received a response, `HeartBeater.Beat(...)` will be called.
//The received response (parameter `response`),
//the time from send to receive (parameter `lastTimeout`),
//the retransmission time before receive (parameter `lastRetryN`)
//and a function `beat` will be the input.
//With those input,
//this method should generate next request (a `protocol.TobeSendRequest`)
//and next send time and retransmission time limit(a `time.Duration` and a `uint64`),
//then call the function `beat(...)`.
//If the function `beat(...)` not called in this methods,
//`Heart.RunBeating(...)` will stop sending and exit.
type HeartBeater interface {
	//In the implementaion of this method `Beat`, you should do the following:
	//
	//1. Cast the `response` into your implementation type (you should implement `protocol.Response` before)
	//
	//2. Generate a `protocol.Request` you want to send to registry, with a `protocol.RequestSendOption`, a timeout limit and a retry limit
	//
	//3. Enclose the `protocol.Request` and `protocol.RequestSendOption` into a `protocol.TobeSendRequest`
	//
	//4. Call the function `beat`, let the input be the generated `protocol.TobeSendRequest`, the timeout limit and the retry limit, or if you want stop the sending, just do not call the `beat` and return.
	Beat(ctx context.Context, response protocol.Response, lastTimeout time.Duration, lastRetryN uint64,
		beat func(protocol.TobeSendRequest, time.Duration, uint64))
}
