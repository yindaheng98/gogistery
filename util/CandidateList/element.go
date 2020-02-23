package CandidateList

import "gogistery/protocol"

type element struct {
	RegistryInfo protocol.RegistryInfo
}

func (e element) GetName() string {
	return e.RegistryInfo.GetRegistryID()
}
