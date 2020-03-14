package registrant

import (
	"context"
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//RegistryCandidateList is the abstract of the candidate registry list.
//This interface is used in Registrant.
type RegistryCandidateList interface {
	//Add the information of a registries to candidate registry list.
	StoreCandidates(ctx context.Context, candidates []protocol.RegistryInfo)

	//Get the information of a candidate registry, except those in "except".
	//If there is no candidate meet the conditions, block until a eligible candidate added.
	GetCandidate(ctx context.Context, except []protocol.RegistryInfo) (protocol.RegistryInfo, time.Duration, uint64)
}
