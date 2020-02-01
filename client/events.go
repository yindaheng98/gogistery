package client

import "gogistery/proto"

type events struct {
	ServerOnline  proto.ServerEmitter //当成功连接到某个服务器时触发事件
	ServerOffline proto.ServerEmitter //当与某个服务器连接断开时触发事件
}
