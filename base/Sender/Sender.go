package Sender

import (
	"errors"
	"gogistery/base"
	"gogistery/base/Errors"
	"gogistery/util/Single"
	"sync"
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

	retryN     uint32      //重试次数
	stopping   bool        //是否准备停止
	stoppingMu *sync.Mutex //启停操作锁
	status     StatusType  //当前的连接状态
	Events     *events     //上/下线事件的处理通过事件触发器完成
	runner     *Single.Thread
}

//新建一个发送端
func New(info base.SenderInfo, proto Protocol, initAddr string, initTimeout time.Duration, initRetryN uint32) *Sender {
	return &Sender{info, proto, nil,
		initAddr, initTimeout, initRetryN, false, new(sync.Mutex),
		STATUS_Disconnected,
		newEvents(),
		Single.NewThread()}
}

//启动发送端轮询协程
func (s *Sender) Connect() {
	s.stoppingMu.Lock()
	defer s.stoppingMu.Unlock()
	s.stopping = false
	go s.runner.Run(s.routine)
}

//发送端轮询协程，程序必须保证此协程任何时候都只有一个在运行
func (s *Sender) routine() {

	/********必要的数据初始化操作********/
	s.Events.Start.Emit() //开始循环则触发启动事件
	s.receiver = nil      //清除之前的连接
	retryN := uint32(0)   //重置重试次数
	disconnected := false //是否在退出前已触发过断连事件
	lastSendTime := time.Now()

	/********轮询退出前要执行的操作********/
	defer func() {
		if !disconnected { //在退出前没有触发过断连事件
			if e := recover(); e != nil {
				s.Events.Disconnect.Emit(Errors.NewLinkError(e.(error), base.NewLinkInfo(s.info, s.receiver)))
			} //那就读取错误触发断连事件
		}
		atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Disconnected)) //进入未连接状态
		s.Events.Stop.Emit()                                                  //退出循环则触发停止事件
	}()

	/********轮询操作********/
	for !s.stopping { //不处于停止状态才继续循环

		/********进行发送和发送时的超时检测********/
		lastSendTime = time.Now()
		protoChan := make(chan ProtoChanElement, 1)
		go s.proto.Send(s.info, s.addr, protoChan) //异步执行发送操作
		go func() {                                //异步执行超时检测函数
			defer func() {
				_ = recover()
			}()
			time.Sleep(s.timeout)                                          //等待一段时间
			protoChan <- ProtoChanElement{nil, errors.New("send timeout")} //发送超时信息
			close(protoChan)                                               //然后关闭通道
		}()

		/********接收并解析发送协议回传的信息********/
		protoInfo, ok := <-protoChan
		if !ok {
			protoInfo = ProtoChanElement{nil, errors.New("no response")}
		}
		receiverInfo := protoInfo.info //拆解发送协议中的信息
		err := protoInfo.error

		/********处理发送协议回传的信息********/
		if err != nil {
			/********如果返回错误就进行如下操作********/
			atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Retrying)) //先进入尝试连接状态
			retryN++                                                          //尝试次数+1
			if retryN <= s.retryN {                                           //如果尝试次数没有超过限制
				s.Events.Retry.Emit(Errors.NewLinkError(err, base.NewLinkInfo(s.info, s.receiver)))
				//就报重试错误
			} else { //如果尝试次数超过了限制
				atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Disconnected)) //那就进入未连接状态
				s.Events.Disconnect.Emit(Errors.NewLinkError(err, base.NewLinkInfo(s.info, s.receiver)))
				disconnected = true //并触发掉线事件
				s.Disconnect()
				break //然后直接退出
			}
		} else if receiverInfo.IsDisconnect() {
			/********不出错但是返回了断开连接的消息那就断开连接********/
			s.Disconnect()
			break
		} else {
			/********不出错就更新地址、延时和重试次数********/
			atomic.StoreUint32((*uint32)(&s.status), uint32(STATUS_Connected)) //进入连接状态
			retryN = uint32(0)                                                 //重置重试次数
			s.addr = receiverInfo.GetAddr()
			s.timeout = receiverInfo.GetTimeout()
			s.retryN = receiverInfo.GetRetryN()
			s.receiver = receiverInfo
			if s.receiver == nil { //如果之前没连过
				s.Events.Connect.Emit(receiverInfo) //就触发上线事件
			} else {
				s.Events.Update.Emit(receiverInfo) //否则触发更新连接事件
			}
			s.receiver = receiverInfo
		}

		/********延时********/
		time.Sleep(s.timeout - time.Now().Sub(lastSendTime))
	}
}

func (s *Sender) Disconnect() {
	s.stoppingMu.Lock()
	defer s.stoppingMu.Unlock()
	if !s.stopping {
		s.proto.SendDisconnect(s.info, s.addr)
	}
	s.stopping = true
}
