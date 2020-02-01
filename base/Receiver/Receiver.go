package Receiver

import (
	"gogistery/base"
	"gogistery/util/Single"
	"gogistery/util/TimeoutMap"
	"sync/atomic"
	"time"
)

//一个Recevier对应进行一种服务的注册，如果需要对多种服务进行注册，创建多个Recevier即可
type Receiver struct {
	proto       Protocol
	info        base.ReceiverInfo
	timeout     time.Duration
	connectionN int
	senders     *TimeoutMap.TimeoutMap

	Events  *events
	runners []*Single.Processor
}

//新建接收器
//
//timeout设置超时时间，threadN设置并行处理的线程数
func New(info base.ReceiverInfo, proto Protocol, timeout time.Duration, threadN int32, connectionN int) *Receiver {
	//构造基础参数
	r := &Receiver{proto, info, timeout, connectionN,
		TimeoutMap.New(), newEvents(), nil}

	//构造线程启停回调
	threadn := threadN
	StartedCallback := func() {
		if atomic.AddInt32(&threadn, -1) <= 0 {
			r.Events.Start.Emit()
		}
	}
	StoppedCallback := func() {
		if atomic.AddInt32(&threadn, 1) >= threadN {
			r.Events.Stop.Emit()
		}
	}

	//构造线程启动器
	runners := make([]*Single.Processor, threadN)
	for i := int32(0); i < threadN; i++ {
		runner := Single.NewProcessor()
		runner.Callback.Started = StartedCallback
		runner.Callback.Stopped = StoppedCallback
		runners[i] = runner
	}

	//返回结果
	r.runners = runners
	return r
}

//获取当前在连的发送端列表
func (r *Receiver) GetSenderInfos() []base.SenderInfo {
	var res []base.SenderInfo
	for _, senderInfo := range r.senders.GetAll() {
		res = append(res, senderInfo.(base.SenderInfo))
	}
	return res
}

func (r *Receiver) Start() {
	for _, runner := range r.runners {
		runner.Start(r.receiverLoop)
	}
}

//接收器循环，同一时刻可能有多个接收器在运行
func (r *Receiver) receiverLoop() {
	senderInfo, err := r.proto.Receive(r.info)
	if err != nil {
		return
	}
	if senderInfo.IsDisconnect() {
		r.senders.Delete(senderInfo.GetID())
		return
	}
	r.senders.UpdateInfo(&element{senderInfo,
		r.Events.Connected,
		r.Events.Disconnected},
		r.timeout)
}

func (r *Receiver) Stop() {
	for _, runner := range r.runners {
		runner.Stop()
	}
}
