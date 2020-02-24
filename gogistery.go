package gogistery

import (
	"github.com/yindaheng98/gogistry/protocol"
	"github.com/yindaheng98/gogistry/registrant"
	"github.com/yindaheng98/gogistry/registry"
)

func NewRegistry(
	Info protocol.RegistryInfo,
	maxRegistrants uint,
	timeoutController registry.TimeoutController,
	ResponseProto protocol.ResponseProtocol) *registry.Registry {
	return registry.New(Info, maxRegistrants, timeoutController, ResponseProto)
}

func NewRegistrant(
	Info protocol.RegistrantInfo,
	regitryN uint,
	CandidateList registrant.RegistryCandidateList,
	retryNController registrant.RetryNController,
	RequestProto protocol.RequestProtocol) *registrant.Registrant {
	return registrant.New(Info, regitryN, CandidateList, retryNController, RequestProto)
}
