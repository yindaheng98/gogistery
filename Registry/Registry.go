package Registry

import (
	"github.com/yindaheng98/go-utility/TimeoutMap"
	"gogistery/Heart"
	"gogistery/Protocol"
	"sync"
	"time"
)

type registrantHandler struct {
	RegistrantInfo Protocol.RegistrantInfo
	registry       *Registry
}

func (info registrantHandler) GetID() string {
	return info.RegistrantInfo.GetRegistrantID()
}

func (info registrantHandler) NewAddedHandler() {
	info.registry.Events.NewConnection.Emit(info.RegistrantInfo)
}
func (info registrantHandler) TimeoutHandler() {
	info.registry.Events.ConnectionTimeout.Emit(info.RegistrantInfo)
}
func (info registrantHandler) DeletedHandler() {
	info.registry.Events.Disconnection.Emit(info.RegistrantInfo)
}

type Registry struct {
	Info      Protocol.RegistryInfo //存储自身信息
	responser *Heart.ResponserHeart //响应器/消息源

	maxRegistrants int                    //最大连接数
	timeoutMap     *TimeoutMap.TimeoutMap //超时计时表
	timeoutMapMu   *sync.RWMutex
	timeoutProto   RegistrantControlProtocol //如何选择timeout

	Events *events
}

func New(Info Protocol.RegistryInfo, maxRegistrants int, timeoutProto RegistrantControlProtocol, sendProto Protocol.ResponseProtocol) *Registry {
	registry := &Registry{
		Info:      Info,
		responser: nil,

		maxRegistrants: maxRegistrants,
		timeoutMap:     TimeoutMap.New(),
		timeoutMapMu:   new(sync.RWMutex),
		timeoutProto:   timeoutProto,

		Events: newEvents(),
	}
	registry.responser = Heart.NewResponserHeart(
		&responserHeartProtocol{registry}, sendProto)
	registry.Events.Error = registry.responser.Event.Error
	return registry
}

//启动直到调用停止才退出
func (r *Registry) Run() {
	r.responser.RunBeating()
}

//停止
func (r *Registry) Stop() {
	r.responser.Stop()
}

//获取当前所有活动连接
func (r *Registry) GetConnections() []Protocol.RegistrantInfo {
	r.timeoutMapMu.RLock()
	defer r.timeoutMapMu.RUnlock()
	infos := make([]Protocol.RegistrantInfo, r.timeoutMap.Count())
	for i, registrant := range r.timeoutMap.GetAll() {
		infos[i] = registrant.(registrantHandler).RegistrantInfo
	}
	return infos
}

//进行一次注册操作，返回指定的下一次心跳的时间限制，如果接受连接则返回true，拒绝连接则返回false
func (r *Registry) register(request Protocol.Request) (time.Duration, uint64, bool) {
	registrantID := request.RegistrantInfo.GetRegistrantID()
	r.timeoutMapMu.Lock()
	defer r.timeoutMapMu.Unlock()
	if request.IsDisconnect() { //如果主动断开连接
		r.timeoutMap.Delete(registrantID) //则直接删除
		return 0, 0, false
	}
	if _, ok := r.timeoutMap.GetElement(registrantID); !ok && r.timeoutMap.Count() >= r.maxRegistrants {
		return 0, 0, false //连接不存在且已达到最大连接数，则拒绝连接
	}
	var timeout time.Duration
	var retryN uint64
	var exists bool
	if _, exists = r.timeoutMap.GetElement(registrantID); !exists {
		timeout, retryN = r.timeoutProto.TimeoutRetryNForNew(request) //不存在则获取新建的timeout
	} else {
		timeout, retryN = r.timeoutProto.TimeoutRetryNForUpdate(request) //存在则获取更新的timeout
	}
	r.timeoutMap.UpdateInfo(
		registrantHandler{request.RegistrantInfo, r}, time.Duration(retryN)*timeout) //否则更新连接
	if exists { //如果存在则说明是更新，触发更新事件
		r.Events.UpdateConnection.Emit(request.RegistrantInfo) //并触发更新事件
	}
	return timeout, retryN, true
}
