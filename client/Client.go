package client

import (
	"gogistery/proto"
	"gogistery/util/SortedSet"
	"log"
)

type client struct {
	selfInfo         proto.ClientInfo            //客户端自己的相关信息
	servers          map[string]proto.ServerInfo //已连接的服务器的相关信息
	candidateServers *SortedSet.SortedSet        //候选服务器的地址
	Events           events                      //用于事件的存储与触发
}

//创建一个客户端
//
//输入info是客户端自身信息，serverN是要同时连接多少个服务器
func Client(info proto.ClientInfo, serverN uint64) *client {
	return &client{info,
		make(map[string]proto.ServerInfo, serverN),
		SortedSet.NewSortedSet(serverN),
		events{proto.NewServerEmitter(), proto.NewServerEmitter()}}
}

//启动此客户端
//
//输入的addr是请求初始化数据的服务器地址
func (cli *client) Start(addr string) error {
	cli.selfInfo.Status = 1
	err := cli.register(addr)
	if err != nil {
		return err
	}
	log.Println("Successfully register to init server: " + addr)
	go cli.poll() //然后启动轮询go程
	log.Println("Gogistery client daemon is started")
	return nil
}

//停止此客户端
func (cli *client) Stop() {
	cli.selfInfo.Status = 0
}

func (cli *client) poll() {
	for {
		if cli.selfInfo.Status != 1 {
			log.Println("Gogistery client daemon is stopped")
			break
		}
		//TODO:每隔一段时间轮询一次
	}
}

func (cli *client) register(addr string) error {
	//更新服务器列表和候选列表的工作在此处完成
	err, server := proto.ClientRegistry(addr, cli.selfInfo)
	if err != nil {
		return err
	}
	for related_server := range server.RelatedServers {
		//这里尽量把新的server放在候选列表的前面，并且不要重复
	}
}
