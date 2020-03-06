package registrant

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

type RegistryCandidateList interface {
	//存入一组候选注册中心
	StoreCandidates(ctx context.Context, candidates []protocol.RegistryInfo)

	//选出一个用于初始化的注册中心信息，并且不能是except中列出的这几个
	GetCandidate(ctx context.Context, except []protocol.RegistryInfo) (protocol.RegistryInfo, time.Duration, uint64)
}
