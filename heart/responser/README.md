# github.com/yindaheng98/gogistry/heart/responser

## Introduction

This package is the "heart" of higher-level protocol in registry. This package has only one exported struct `Heart`, with only one exported method:

```go
func (h *Heart) RunBeating(ctx context.Context)
```


1. Call the method `protocol.ResponseProtocol.Response` you have implemented to receive a request and a function for sending back response
2. Call the method `HeartBeater.Beat` you have implemented to get a response
3. Send back the response using the received sending-back function
4. Loop step 1~3 until a stop signal arrived

When called, this function will block your goroutine until returned.

## Interface

The only interface you should implment in this package is the bridges for higher-level protocol `HeartBeater`.

### `HeartBeater`

`HeartBeater` is the controller of the `Heart`, it has only one method:

```go
type HeartBeater interface {
	Beat(ctx context.Context, request protocol.Request) protocol.TobeSendResponse
}
```

The method was designed to generate a response (`protocol.Response` and `protocol.ResponseSendOption` in `protocol.TobeSendResponse`) from a received request (`Protocol.Request`).

In the implementaion of this method, you should do the following:

1. Cast the `request` into your implementaion type
2. Generate a `protocol.Response` you want to send back to the request and the send option `protocol.ResponseSendOption`
3. Enclose the `protocol.Response` and `protocol.ResponseSendOption` into a `protocol.TobeSendResponse`
4. Return the `protocol.TobeSendResponse`
