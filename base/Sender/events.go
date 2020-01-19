package Sender

import (
	"gogistery/base/Emitter"
)

type events struct {
	Connect    *Emitter.ReceiverInfoEmitter //连接成功
	Retry      *Emitter.LinkErrorEmitter    //重试
	Disconnect *Emitter.LinkErrorEmitter    //断开
	Error      *Emitter.LinkErrorEmitter    //连接失败
	Start      *Emitter.EmptyEmitter        //启动
	Stop       *Emitter.EmptyEmitter        //停止
}

func newEvents() *events {
	return &events{Emitter.NewReceiverInfoEmitter(),
		Emitter.NewLinkErrorEmitter(),
		Emitter.NewLinkErrorEmitter(),
		Emitter.NewLinkErrorEmitter(),
		Emitter.NewEmptyEmitter(),
		Emitter.NewEmptyEmitter()}
}

func (e *events) EnableAll() {
	e.Connect.Enable()
	e.Retry.Enable()
	e.Disconnect.Enable()
	e.Error.Enable()
	e.Start.Enable()
	e.Stop.Enable()
}

func (e *events) DisableAll() {
	e.Connect.Disable()
	e.Retry.Disable()
	e.Disconnect.Disable()
	e.Error.Disable()
	e.Start.Disable()
	e.Stop.Disable()
}
