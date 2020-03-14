package TimeoutController

import (
	"github.com/yindaheng98/gogistry/protocol"
	"math"
	"time"
)

//LogTimeoutController is a simple implementation of TimeoutController.
type LogTimeoutController struct {
	minT time.Duration //最小Timeout
	maxT time.Duration //最大Timeout
	cT   float64       //从最小到最大的增长系数
	tMap map[string]time.Duration
}

//NewLogTimeoutController returns the pointer to a LogTimeoutController.
func NewLogTimeoutController(minT time.Duration, maxT time.Duration, cT float64) *LogTimeoutController {
	return &LogTimeoutController{minT, maxT, cT,
		make(map[string]time.Duration)}
}

//timeout(0)=minT
func (p LogTimeoutController) TimeoutForNew(request protocol.Request) time.Duration {
	p.tMap[request.RegistrantInfo.GetRegistrantID()] = p.minT
	return p.minT

}

//timeout(n)=timeout(n-1)+(maxT-timeout(n-1))/cT
func (p LogTimeoutController) TimeoutForUpdate(request protocol.Request) time.Duration {
	t := p.tMap[request.RegistrantInfo.GetRegistrantID()]
	t += time.Duration(math.Floor(float64(p.maxT-t) / p.cT))
	p.tMap[request.RegistrantInfo.GetRegistrantID()] = t
	return t
}
