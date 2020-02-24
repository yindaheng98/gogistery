package registry

import "github.com/yindaheng98/gogistry/protocol"

type registrantTimeoutType struct {
	RegistrantInfo protocol.RegistrantInfo
	events         *events
}

func (r registrantTimeoutType) GetID() string {
	return r.RegistrantInfo.GetRegistrantID()
}

func (r registrantTimeoutType) NewAddedHandler() {
	r.events.NewConnection.Emit(r.RegistrantInfo)
}
func (r registrantTimeoutType) UpdatedHandler() {
	r.events.UpdateConnection.Emit(r.RegistrantInfo)
}
func (r registrantTimeoutType) TimeoutedHandler() {
	r.events.ConnectionTimeout.Emit(r.RegistrantInfo)
}
func (r registrantTimeoutType) DeletedHandler() {
	r.events.Disconnection.Emit(r.RegistrantInfo)
}
