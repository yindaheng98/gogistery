package Errors

import "gogistery/base"

//用于指示在“从哪个sender到哪个receiver的发送过程中发生了什么error”
type Error struct {
	error
	receiver base.ReceiverInfo
	sender   base.SenderInfo
}

func New(err error, receiver base.ReceiverInfo, sender base.SenderInfo) Error {
	return Error{error: err, sender: sender, receiver: receiver}
}

func (e *Error) Receiver() base.ReceiverInfo {
	return e.receiver
}
func (e *Error) Sender() base.SenderInfo {
	return e.sender
}
