package TimeoutController

import (
	"github.com/yindaheng98/gogistry/protocol"
	"math"
	"time"
)

//LogTimeoutController is a simple implementation of TimeoutController.
type LogTimeoutController struct {
	MinimumTime    time.Duration //最小Timeout
	MaximumTime    time.Duration //最大Timeout
	IncreaseFactor float64       //从最小到最大的增长系数
	tMap           map[string]time.Duration
}

//DefaultLogTimeoutController returns the pointer to a LogTimeoutController with default member value.
func DefaultLogTimeoutController() *LogTimeoutController {
	return &LogTimeoutController{1e9, 10e9, 2,
		make(map[string]time.Duration)}
}

//Timeout(0)=MinimumTime
func (p LogTimeoutController) TimeoutForNew(request protocol.Request) time.Duration {
	p.tMap[request.RegistrantInfo.GetRegistrantID()] = p.MinimumTime
	return p.MinimumTime

}

//Timeout(n)=Timeout(n-1)+(MaximumTime-Timeout(n-1))/IncreaseFactor
func (p LogTimeoutController) TimeoutForUpdate(request protocol.Request) time.Duration {
	t := p.tMap[request.RegistrantInfo.GetRegistrantID()]
	t += time.Duration(math.Floor(float64(p.MaximumTime-t) / p.IncreaseFactor))
	p.tMap[request.RegistrantInfo.GetRegistrantID()] = t
	return t
}
