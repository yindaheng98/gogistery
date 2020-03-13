package protocol

import (
	"context"
	"fmt"
)

//ReceivedResponse is used to store the response and error received by registrant.
type ReceivedResponse struct {
	Response Response
	Error    error
}

//ReceivedRequest is used to store the request and error received by registry.
type ReceivedRequest struct {
	Request Request
	Error   error
}

//TobeSendRequest is used to store the request that is going to be sent by registrant to registry.
//Option is the option information for sending (encoding, encryption, etc).
type TobeSendRequest struct {
	Request Request
	Option  RequestSendOption
}

func (r TobeSendRequest) String() string {
	return fmt.Sprintf("TobeSendRequest{Request:%s,Option:%s}", r.Request.String(), r.Option.String())
}

//TobeSendResponse is used to store the response that is going to be sent back by registry to registrant.
//Option is the option information for sending (encoding, encryption, etc).
type TobeSendResponse struct {
	Response Response
	Option   ResponseSendOption
}

func (r TobeSendResponse) String() string {
	return fmt.Sprintf("TobeSendResponse{Response:%s,Option:%s}", r.Response.String(), r.Option.String())
}

//RequestProtocol defines how a request should be sent from registrant to registry (by http, grpc, etc).
//It should be implement by user.
//
//Higher-level protocol (`"heart/requester".Heart`) will call this function to send a request according to some option (they will be a `Request` and a `RequestSendOption`, enclosed into `TobeSendRequest`) via lower-level protocol (i.e. http, grpc, etc, implemented by you), and receive the response.
type RequestProtocol interface {
	//In the implementation of this method, you should:
	//
	//1. Get a `TobeSendRequest` from the `requestChan`
	//
	//2. Get a `Request` and a `RequestSendOption` from `TobeSendRequest`
	//
	//3. Send the `Request` according to `RequestSendOption` via your protocol
	//
	//4. Waiting for a response in your protocol, and enclose the response into a `Response`, if there is an error occured, you should generate an `error`
	//
	//5. Enclose the `Response` and the `error` (if exists) into a `ReceivedResponse`
	//
	//6. Put the `ReceivedResponse` into the `responseChan`
	Request(ctx context.Context, requestChan <-chan TobeSendRequest, responseChan chan<- ReceivedResponse)
}

//ResponseProtocol defines how a response should be sent back from registrant to registry (by http, grpc, etc).
//It should be implement by user.
//
//Higher-level protocol(`"heart/responser".Heart`) will call this function in a loop, to receive the request (`Request`) from lower-level protocol (i.e. http, grpc, etc, implemented by you) from registrant and send back the generated response and sending option (`Response` and `ResponseSendOption`).
type ResponseProtocol interface {
	//In the implementation of this method, you should:
	//
	//1. Waiting for a request in your protocol, if there is an error occured, you should generate an `error`
	//
	//2. Enclose the request into a `Request`
	//
	//3. Enclose the `Request` and the `error` (if exists) into a `ReceivedRequest`
	//
	//4. Put the `ReceivedRequest` into `requestChan`
	//
	//5. Waiting for `responseChan` to send back a `TobeSendResponse`
	//
	//6. Get a `Response` and a `ResponseSendOption` from the `TobeSendResponse`
	//
	//7. Send back the `Response` according to the `ResponseSendOption` via your protocol
	Response(ctx context.Context, requestChan chan<- ReceivedRequest, responseChan <-chan TobeSendResponse)
}
