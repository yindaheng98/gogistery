# Heart

## Introduction

This package is the "heart" of higher-level protocol in gogistry (lower-level protocol should be implemented by yourself in the package `Protocol`). It contains the only way to run gogistry system (the only way to "beat" the "heart"). All the other functions or interfaces in gogistry system will be called in this package.

This package was designed for the controlling of the "heartbeat" sequence in gogistry. The "heartbeat" means a list of registrant is keeping in each registry, and can be accessed from some interface; registrant should frequently send messages to registry, or will be deleted from the registrant list. This package is going to control the frequency of "heartbeat" sending from registrant, and decide when should a registry delete a registrant from the registrant list.

A usage example is in [example/Heart](example/Heart)

## You should implement the following interfaces

The interfaces you should implment in this package are the bridges for higher-level protocol.

### `RequestSendOption`

`RequestSendOption` is the option interface for `RequesterHeart` it designed to contain necessary information for sending a `Protocol.Request` in `RequesterHeart`.

This interface has three methods:

```go
type RequestSendOption interface {
	GetTimeout() time.Duration
	GetRetryN() uint64
	String() string
}
```

When sending every `Request`, `RequesterHeart` will attempt for `RequestSendOption.GetRetryN()` times, each time waiting the response during `RequestSendOption.GetTimeout()`, if response was not arrived in this duration, an error "timeout" will be throw.

### `RequesterHeartProtocol`

`RequesterHeartProtocol` will run in registrant, it has only one method:

```go
type RequesterHeartProtocol interface {
	Beat(response Protocol.Response, beat func(Protocol.TobeSendRequest))
}
```

The method was designed to get a response (`Protocol.Response`) from the registry, then generate next request (`Protocol.Request` and `RequestSendOption`).

In the implementaion of this method, you should do the following:

1. Cast the `response` into your implementation type (you should implement `Protocol.Response` before)
2. Generate a `Protocol.Request` you want to send to registry, and put the timeout limit with retry limit into a `RequestSendOption`
3. Enclose the `Request` and `RequestSendOption` into a `Protocol.TobeSendRequest`
4. Call the function `beat`, let the input be the generated `Protocol.TobeSendRequest`

### ResponserHeartProtocol

`ResponserHeartProtocol` will run in registry, it has only one method:

```go
type ResponserHeartProtocol interface {
	Beat(request Protocol.Request) Protocol.TobeSendResponse
}
```

The method was designed to generate a response (`Protocol.Response` and `Protocol.ResponseSendOption`) from a received request (`Protocol.Request`).

In the implementaion of this method, you should do the following:

1. Cast the `request` into your implementaion type
2. Generate a `Protocol.Response` you want to send back to the request and the send option `Protocol.ResponseSendOption`
3. Enclose the `Protocol.Response` and `Protocol.ResponseSendOption` into a `Protocol.TobeSendResponse`
4. Return the `Protocol.TobeSendResponse`

## How to use the package

### function `NewRequesterHeart` and struct `RequesterHeart`

Struct `RequesterHeart`, as its name, is designed for request. It is the "heart" of registrant. It will send "beat" to registry in a loop, until the registry reject it or occured some error.

Function `NewRequesterHeart`, as its name, to generate a new instance of `RequesterHeart`:

```go
func NewRequesterHeart(heartProto RequesterHeartProtocol, beatProto Protocol.RequestProtocol) *RequesterHeart
```

Its input is an instance of `RequesterHeartProtocol` and an instance of `Protocol.RequestProtocol`(you should have implemented those two interfaces before). Its output is a pointer of `RequesterHeart`, which will act according to your protocol instances.

### Methods of struct `RequesterHeart`

`RequesterHeart` has only one method:

```go
func (h *RequesterHeart) RunBeating(initRequest Protocol.TobeSendRequest) error
```

This methods will do following things:

1. Call the method `Protocol.RequestProtocol.Send` you have implemented to send `initRequest` to registry, and get a response
2. Call the method `RequesterHeartProtocol.Beat` to get next request
3. Call the method `Protocol.RequestProtocol.Send` you have implemented to send the request to registry, and get the other response
4. Loop step 2 and step 3 until some error occured or function `beat` was not called in `RequesterHeartProtocol.Beat`

When called, this function will block your goroutine until returned.

### function `NewResponserHeart` and struct `ResponserHeart`

Struct `ResponserHeart`, as its name, is designed for response. It is the "heart" of registry. It will receive "beat" from registrant in a loop, then send back a response.

Function `NewResponserHeart`, as its name, to generate a new instance of `ResponserHeart`:

```go
func NewResponserHeart(heartProto ResponserHeartProtocol, beatProto Protocol.ResponseProtocol) *ResponserHeart
```

The input is an instance of `ResponserHeartProtocol` and an instance of `Protocol.ResponseProtocol`(you should have implemented those two interfaces before). Its output is a pointer of `ResponserHeart`, which will act according to your protocol instances.

### Methods of struct `ResponserHeart`

`ResponserHeart` has only two method. One is:

```go
func (h *ResponserHeart) RunBeating()
```

This methods will do following things:

1. Call the method `Protocol.ResponseProtocol.Recv` you have implemented to receive a request and a function for sending back response
2. Call the method `RequesterHeartProtocol.Beat` you have implemented to get a response
3. Send back the response using the received sending-back function
4. Loop step 1~3 until a stop signal arrived

When called, this function will block your goroutine until returned.

The other method of `ResponserHeart` is:

```go
func (h *ResponserHeart) Stop()
```

This method is just send a stop signal to `ResponserHeart.RunBeating` and waiting for its return. When called, this function will block your goroutine until `ResponserHeart.RunBeating` returned.