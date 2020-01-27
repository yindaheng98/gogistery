package TimeoutMap

import (
	"testing"
	"time"
)

type TestElement struct {
	id string
}

func (t TestElement) GetID() string {
	return t.id
}

func TestTimeoutMap(t *testing.T) {
	tm := New(1e9, 10)
	tm.UpdateInfo(TestElement{"test1"})
	tm.UpdateInfo(TestElement{"test2"})
	tm.UpdateInfo(TestElement{"test3"})
	t.Log(tm.GetAll())
	tm.Delete("test2")
	t.Log(tm.GetAll())
	go tm.UpdateInfo(TestElement{"test4"})
	time.Sleep(2e8)
	go tm.UpdateInfo(TestElement{"test5"})
	go tm.UpdateInfo(TestElement{"test6"})
	time.Sleep(3e8)
	t.Log(tm.GetAll())
	time.Sleep(5e8)
	t.Log(tm.GetAll())

}
