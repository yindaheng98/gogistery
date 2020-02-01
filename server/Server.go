package server

import "gogistery/proto"

type Server struct {
	ServiceType   string //服务类型，客户端发来的请求必须与此相匹配
}

func New() *Server {
	return &Server{"s0"}
}

func (s *Server) Start(where string) {

}

func(s*Server)On(event string,info proto.ClientInfo){

}