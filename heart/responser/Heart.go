package responser

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
)

//Heart can receive the requests and send back response in a loop,
//until stopped by context.
type Heart struct {
	beater    HeartBeater
	responser *responser

	//This function will be called when an error occurred.
	ErrorHandler func(error)
}

//NewHeart returns the pointer to a heart.
func NewHeart(beater HeartBeater, ResponseProto protocol.ResponseProtocol) *Heart {
	return &Heart{beater,
		newResponser(ResponseProto),
		func(error) {}}
}

//RunBeating is a loop to receive the requests and send back response.
//It will do the following:
//
//1. Call the method `protocol.ResponseProtocol.Response` you have implemented to receive a request and a function for sending back response
//
//2. Call the method `HeartBeater.Beat` you have implemented to get a response
//
//3. Send back the response using the received sending-back function
//
//4. Loop step 1~3 until a stop signal arrived
//
//When called, this function will block your goroutine until returned.
func (h *Heart) RunBeating(ctx context.Context) {
	for {
		var request protocol.Request
		var err error
		var responseFunc func(protocol.TobeSendResponse)
		okChan := make(chan bool, 1)
		go func() {
			request, err, responseFunc = h.responser.Recv(ctx)
			okChan <- true
		}()
		select {
		case <-okChan:
			if err != nil {
				h.ErrorHandler(err)
			} else {
				responseFunc(h.beater.Beat(ctx, request))
			}
		case <-ctx.Done():
			return
		}
	}
}
