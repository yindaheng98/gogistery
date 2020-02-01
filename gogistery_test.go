package gogistery

import (
	"gogistery/proto"
	"testing"
)

func TestServerInfo(t *testing.T) {
	t.Log("服务器和客户端信息测试")
	sInfo := proto.ServerInfo{
		ServiceType:    "testService",
		ID:             "test",
		Addr:           "test.test.wxstc",
		RelatedServers: []string{"a.test.o"},
	}
	sInfoB, err := sInfo.String()
	if err != nil {
		t.Error(err)
	}
	sInfoS := string(sInfoB)
	t.Log(sInfoS)

	sInfoP, err := proto.ParseServer([]byte(sInfoS))
	if err != nil {
		t.Error(err)
	}
	t.Log(sInfoP)
}
