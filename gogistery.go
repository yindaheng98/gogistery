package gogistery

import (
	"gogistery/protocol"
	"gogistery/registrant"
	"gogistery/registry"
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
