# github.com/yindaheng98/gogistry/heart/requester

## Introduction

This package is the "heart" of higher-level protocol in registrant. This package has only one exported struct `Heart`, with only one exported method:

```go
func (h *Heart) RunBeating(ctx context.Context, initRequest protocol.TobeSendRequest, initTimeout time.Duration, initRetryN uint64) error
```

This methods will send "beat" to registry in a loop, until the registry reject it or occured some error. It will do following things:

1. Call the method `protocol.RequestProtocol.Request` you have implemented to send `initRequest` to registry, and get a response
2. Call the method `HeartBeater.Beat` to get next request
3. Call the method `protocol.RequestProtocol.Request` you have implemented to send the request to registry, and get the other response
4. Loop step 2 and step 3 until some error occured or function `beat` was not called in `HeartBeater.Beat`

When called, this function will block your goroutine until returned.

## Interface

The only interface you should implment in this package is the bridges for higher-level protocol `HeartBeater`.

### `HeartBeater`

`HeartBeater` is the controller of the `Heart`, it has only one method:

```go
type HeartBeater interface {
	Beat(ctx context.Context, response protocol.Response, lastTimeout time.Duration, lastRetryN uint64, beat func(protocol.TobeSendRequest, time.Duration, uint64))
}
```

In `Heart.RunBeating(...)`, when sent a request and received a response, this method will be called. The received response (parameter `response`), the time from send to receive (parameter `lastTimeout`), the retransmission time before receive (parameter `lastRetryN`) and a function `beat` will be the input. With those input, this method should generate next request (a `protocol.TobeSendRequest`) and next send time and retransmission time limit(a `time.Duration` and a `uint64`), then call the function `beat(...)`. If the function `beat(...)` not called in this methods, `Heart.RunBeating(...)` will stop sending and exit.

So in the implementaion of this method, you should do the following:

1. Cast the `response` into your implementation type (you should implement `protocol.Response` before)
2. Generate a `protocol.Request` you want to send to registry, with a `protocol.RequestSendOption`, a timeout limit and a retry limit
3. Enclose the `protocol.Request` and `protocol.RequestSendOption` into a `protocol.TobeSendRequest`
4. Call the function `beat`, let the input be the generated `protocol.TobeSendRequest`, the timeout limit and the retry limit, or if you want stop the sending, just do not call the `beat` and return.