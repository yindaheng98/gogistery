package gogistery

import (
	"gogistery/InfoStructs"
	"testing"
)

func TestServerInfo(t *testing.T) {
	t.Log("服务器和客户端信息测试")
	sInfo := InfoStructs.ServerInfo{
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

	sInfoP, err := InfoStructs.ParseServer([]byte(sInfoS))
	if err != nil {
		t.Error(err)
	}
	t.Log(sInfoP)
}
