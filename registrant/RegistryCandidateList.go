package registrant

import (
	"gogistery/protocol"
	"time"
)

type RegistryCandidateList interface {
	//存入一组候选注册中心
	StoreCandidates(response protocol.Response)

	//选出一个用于初始化的注册中心信息，并且不能是except中列出的这几个
	GetCandidate(except []protocol.RegistryInfo) (protocol.RegistryInfo, time.Duration, uint64)
}
