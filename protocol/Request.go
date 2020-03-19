package protocol

import "fmt"

//ResponseSendOption is the option information for response sending (encoding, encryption, etc).
type ResponseSendOption interface {
	String() string
}

//RegistrantInfo contains the information for registrant.
//It will be send from registrant to registry within the request.
//It should be implement by user.
type RegistrantInfo interface {

	//Returns the unique ID of the registrant
	GetRegistrantID() string

	//Returns the type of the service
	GetServiceType() string

	//Returns the option when the registry send back the response
	GetResponseSendOption() ResponseSendOption

	String() string
}

//Request is the request that registrant send to registry.
//It contains the information for registrant "RegistrantInfo" and a connection flag "Disconnect".
type Request struct { //服务器端收到的请求
	RegistrantInfo RegistrantInfo
	Disconnect     bool
}

//Get the value of connection flag "Disconnect".
func (r Request) IsDisconnect() bool {
	return r.Disconnect
}

func (r Request) String() string {
	return fmt.Sprintf(`{"type":"github.com/yindaheng98/gogistry/protocol.Request",
	"RegistrantInfo":%s,"Disconnect":"%t"}`, r.RegistrantInfo.String(), r.Disconnect)
}
