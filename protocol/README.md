# Protocol

## Introduction

This package is the basic of the gogistry. Before using gogistry, you should implement the interfaces using the transmission protocol you like (i.e. http, grpc, etc) in this package first.

## How to implement

There are 6 interface you should implement: `RequestSendOption`, `ResponseSendOption`, `RegistrantInfo`, `RegistryInfo`, `RequestProtocol` and `ResponseProtocol`.

### Implement `RequestSendOption` and `ResponseSendOption`

This two interface is designed for carry on those informations which is related to the message sending or receiving .

### Implement the interface `RequestProtocol`

`RequestProtocol` will running in registrants. It has just one method:

```go
type RequestProtocol interface {
	Request(requestChan <-chan TobeSendRequest, responseChan chan<- ReceivedResponse)
}
```

Higher-level protocol (`"heart/requester".Heart`) will call this function to send a request according to some option (they will be a `Request` and a `RequestSendOption` you just implemented, enclosed into `TobeSendRequest`) via lower-level protocol (i.e. http, grpc, etc, implemented by you), and receive the response.

In the implementation of `RequestProtocol`, you should do the following actions in `RequestProtocol.Request`:

1. Get a `TobeSendRequest` from the `requestChan`
2. Get a `Request` and a `RequestSendOption` from `TobeSendRequest`
3. Send the `Request` according to `RequestSendOption` via your protocol
4. Waiting for a response in your protocol, and enclose the response into a `Response`, if there is an error occured, you should generate an `error`
5. Enclose the `Response` and the `error` (if exists) into a `ReceivedResponse`
6. Put the `ReceivedResponse` into the `responseChan`

### Implement the interface `ResponseProtocol`

`ResponseProtocol` will running in registry. It has just one method: 

```go
type ResponseProtocol interface {
	Response(requestChan chan<- ReceivedRequest, responseChan <-chan TobeSendResponse)
}
```

Higher-level protocol(`"heart/responser".Heart`) will call this function in a loop, to receive the request (`Request`) from lower-level protocol (i.e. http, grpc, etc, implemented by you) from registrant and send back the generated response and sending option (`Response` and `ResponseSendOption`).

In the implementation of `ResponseProtocol`, you should do the following actions in `ResponseProtocol.Response`:

1. Waiting for a request in your protocol, if there is an error occured, you should generate an `error`
2. Enclose the request into a `Request`
3. Enclose the `Request` and the `error` (if exists) into a `ReceivedRequest`
4. Put the `ReceivedRequest` into `requestChan`
5. Waiting for `responseChan` to send back a `TobeSendResponse`
6. Get a `Response` and a `ResponseSendOption` from the `TobeSendResponse`
7. Send back the `Response` according to the `ResponseSendOption` via your protocol
