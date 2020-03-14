package CandidateList

import "github.com/yindaheng98/gogistry/protocol"

type element struct {
	RegistryInfo protocol.RegistryInfo
}

func (e element) GetName() string {
	return e.RegistryInfo.GetRegistryID()
}
