package registry

import (
	"context"
	"github.com/yindaheng98/go-utility/TimeoutMap"
	"github.com/yindaheng98/gogistry/heart/responser"
	"github.com/yindaheng98/gogistry/protocol"
	"sync"
	"time"
)

//Registrant stands for the "registry" in gogistry.
//A "registry" will receive requests from registrant, record them and send back response in a loop.
type Registry struct {

	//Info contains the information of this Registry.
	Info  protocol.RegistryInfo //存储自身信息
	heart *responser.Heart      //响应器/消息源

	maxRegistrants    uint64                 //最大连接数
	timeoutMap        *TimeoutMap.TimeoutMap //超时计时表
	timeoutMapMu      *sync.RWMutex
	timeoutController TimeoutController //如何选择timeout

	//Events contains 5 emitters to record running events.
	Events *events
}

//New returns the pointer to a Registry.
func New(Info protocol.RegistryInfo, maxRegistrants uint64, timeoutController TimeoutController, ResponseProto protocol.ResponseProtocol) *Registry {
	registry := &Registry{
		Info:  Info,
		heart: nil,

		maxRegistrants:    maxRegistrants,
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

//Run will start the loop of request receive and response send.
func (r *Registry) Run(ctx context.Context) {
	r.heart.RunBeating(ctx)
}

//Registrant maintains a connection list,
//which can be accessed through GetConnections.
//This connection list contains the information (`protocol.RegistrantInfo`) of the registrants connecting with the Registry.
func (r *Registry) GetConnections() []protocol.RegistrantInfo {
	r.timeoutMapMu.RLock()
	defer r.timeoutMapMu.RUnlock()
	timeoutMapEls := r.timeoutMap.GetAll()
	infos := make([]protocol.RegistrantInfo, len(timeoutMapEls))
	for i, registrant := range timeoutMapEls {
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
		return 3600 * 24 * 1e9, false     //直接拒绝连接
	}
	var timeout time.Duration
	var exists bool
	if _, exists = r.timeoutMap.GetElement(registrantID); exists { //如果连接已存在
		timeout = r.timeoutController.TimeoutForUpdate(request) //则获取更新的timeout
	} else if uint64(r.timeoutMap.Count()) < r.maxRegistrants { //如果连接不存在但可以接受新连接
		timeout = r.timeoutController.TimeoutForNew(request) //则获取新建的timeout
	} else { //如果连接不存在且又不能接受新连接
		return 3600 * 24 * 1e9, false //拒绝，让对方明天再来
	}
	r.timeoutMap.UpdateInfo(
		registrantTimeoutType{request.RegistrantInfo, r.Events}, timeout) //更新连接
	if exists { //如果存在则说明是更新
		r.Events.UpdateConnection.Emit(request.RegistrantInfo) //触发更新事件
	}
	return timeout, true
}
