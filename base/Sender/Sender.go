package Sender

import (
	"gogistery/base"
	"gogistery/base/Errors"
	"sync/atomic"
	"time"
)

//发送器类，一个发送器只负责对一个接收器发信息；如果需要同时对多个接收器发送，创建多个发送器即可
//
//服务端发回备用服务器的操作在更高一层实现
type Sender struct {
	info     base.SenderInfo   //存储自身信息
	proto    Protocol          //存储使用的协议
	receiver base.ReceiverInfo //存储连接的服务器的信息
	addr     string            //要向何处发送
	timeout  time.Duration     //超时时间

	retryN   uint32     //重试次数
	started  uint32     //是否已经开始运行（用于控制goroutine）
	stopping bool       //是否准备停止
	status   StatusType //当前的连接状态
	Events   *events    //上/下线事件的处理通过事件触发器完成
}

//新建一个发送端
func New(info base.SenderInfo, proto Protocol, initAddr string, initTimeout time.Duration, initRetryN uint32) *Sender {
	return &Sender{info, proto, nil,
		initAddr, initTimeout, initRetryN, 0, false,
		STATUS_Disconnected,
		newEvents()}
}

//启动发送端轮询协程
func (s *Sender) Connect() {
	s.stopping = false
	go s.routine()
}

//发送端轮询协程，程序必须保证此协程任何时候都只有一个在运行
func (s *Sender) routine() {
	if !atomic.CompareAndSwapUint32(&s.started, 0, 1) { //处于停止状态才启动
		return
	}
	defer func() { //routine退出时也要修改状态
		if atomic.CompareAndSwapUint32(&s.started, 1, 0) {
			atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Disconnected)) //进入未连接状态
			s.Events.Stop.Emit()                                                  //退出循环则触发停止事件
		}
	}()
	s.Events.Start.Emit() //开始循环则触发启动事件
	s.receiver = nil      //清除之前的连接
	retryN := uint32(0)   //重置重试次数
	for !s.stopping {     //不处于停止状态才继续循环
		receiverInfo, err := s.proto.Send(s.info, s.addr, s.timeout) //执行发送操作
		if err != nil {                                              //如果出错
			atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Retrying)) //先进入尝试连接状态
			retryN++                                                          //尝试次数+1
			if retryN <= s.retryN {                                           //如果尝试次数没有超过限制
				s.Events.Retry.Emit(Errors.NewLinkError(err, base.NewLinkInfo(s.info, s.receiver)))
				//就报重试错误
			} else { //如果尝试次数超过了限制
				atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Disconnected)) //那就进入未连接状态
				s.Events.Disconnect.Emit(Errors.NewLinkError(err, base.NewLinkInfo(s.info, s.receiver)))
				//并触发掉线事件
				s.Disconnect()
				break //然后直接退出
			}
		} else { //不出错就更新地址、延时和重试次数
			atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Connected)) //进入连接状态
			retryN = uint32(0)                                                 //重置重试次数
			s.addr = receiverInfo.GetAddr()
			s.timeout = receiverInfo.GetTimeout()
			s.retryN = receiverInfo.GetRetryN()
			if s.receiver == nil { //如果之前没连过
				s.receiver = receiverInfo
				s.Events.Connect.Emit(s.receiver) //就触发上线事件
			}
			s.receiver = receiverInfo
		}
		time.Sleep(s.timeout) //然后延时继续
	}
}

func (s *Sender) Disconnect() {
	s.stopping = true
}
