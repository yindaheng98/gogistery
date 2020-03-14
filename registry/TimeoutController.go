package registry

import (
	"github.com/yindaheng98/gogistry/protocol"
	"time"
)

//TimeoutController is used to control the registrant timer in registry.
//Each registry maintains a TimeoutMap (https://godoc.org/github.com/yindaheng98/go-utility/TimeoutMap).
//Every connected registrant has a timer in the TimeoutMap.
//When received a request from a registrant, the timer of the registrant in TimeoutMap will be reset.
//If the timer of a registrant is not reset for too long, the registrant will be regarded as disconnected and be deleted.
type TimeoutController interface {

	//If received a request from a registrant not exists in TimeoutMap,
	//this function will be called to determine the time of the timer.
	TimeoutForNew(request protocol.Request) time.Duration

	//If received a request from a registrant exists in TimeoutMap,
	//this function will be called to determine the time of the timer before reset.
	TimeoutForUpdate(request protocol.Request) time.Duration
}
