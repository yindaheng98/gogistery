package registry

import (
	"github.com/yindaheng98/go-utility/TimeoutMap"
	"gogistery/heart/responser"
	"gogistery/protocol"
	"sync"
	"time"
)

type Registry struct {
	Info  protocol.RegistryInfo //存储自身信息
	heart *responser.Heart      //响应器/消息源

	maxRegistrants    int                    //最大连接数
	timeoutMap        *TimeoutMap.TimeoutMap //超时计时表
	timeoutMapMu      *sync.RWMutex
	timeoutController TimeoutController //如何选择timeout

	Events *events
}

func New(Info protocol.RegistryInfo, maxRegistrants uint, timeoutController TimeoutController, ResponseProto protocol.ResponseProtocol) *Registry {
	registry := &Registry{
		Info:  Info,
		heart: nil,

		maxRegistrants:    int(maxRegistrants),
		timeoutMap:        TimeoutMap.New(),
		timeoutMapMu:      new(sync.RWMutex),
		timeoutController: timeoutController,

		Events: newEvents(),
	}
	registry.heart = responser.NewHeart(&beater{registry}, ResponseProto)
	registry.heart.ErrorHandler = func(err error) {
		registry.Events.Error.Emit(err)
	}
	return registry
}

//启动直到调用停止才退出
func (r *Registry) Run() {
	r.heart.RunBeating()
}

//停止
func (r *Registry) Stop() {
	r.heart.Stop()
}

//获取当前所有活动连接
func (r *Registry) GetConnections() []protocol.RegistrantInfo {
	r.timeoutMapMu.RLock()
	defer r.timeoutMapMu.RUnlock()
	infos := make([]protocol.RegistrantInfo, r.timeoutMap.Count())
	for i, registrant := range r.timeoutMap.GetAll() {
		infos[i] = registrant.(registrantTimeoutType).RegistrantInfo
	}
	return infos
}

//进行一次注册操作，返回指定的下一次心跳的时间限制，如果接受连接则返回true，拒绝连接则返回false
func (r *Registry) register(request protocol.Request) (time.Duration, bool) {
	registrantID := request.RegistrantInfo.GetRegistrantID()
	r.timeoutMapMu.Lock()
	defer r.timeoutMapMu.Unlock()
	if request.IsDisconnect() { //如果主动断开连接
		r.timeoutMap.Delete(registrantID) //则直接删除
		return 0, false
	}
	var timeout time.Duration
	var exists bool
	if _, exists = r.timeoutMap.GetElement(registrantID); exists { //如果连接已存在
		timeout = r.timeoutController.TimeoutForUpdate(request) //则获取更新的timeout
	} else if r.timeoutMap.Count() < r.maxRegistrants { //如果连接不存在但可以接受新连接
		timeout = r.timeoutController.TimeoutForNew(request) //则获取新建的timeout
	} else { //如果连接不存在且又不能接受新连接
		return 0, false //直接拒绝连接
	}
	r.timeoutMap.UpdateInfo(
		registrantTimeoutType{request.RegistrantInfo, r.Events}, timeout) //更新连接
	if exists { //如果存在则说明是更新
		r.Events.UpdateConnection.Emit(request.RegistrantInfo) //触发更新事件
	}
	return timeout, true
}
