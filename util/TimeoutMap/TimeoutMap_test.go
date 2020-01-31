package TimeoutMap

import (
	"fmt"
	"testing"
	"time"
)

type TestElement struct {
	id string
}

func (e TestElement) GetID() string {
	return e.id
}

func (e TestElement) NewAddedHandler() {
	fmt.Printf("Element %s was new added.\n", e.id)
}

func (e TestElement) TimeoutHandler() {
	fmt.Printf("Element %s is timeout.\n", e.id)
}

func TestTimeoutMap(t *testing.T) {
	tm := New()
	tm.UpdateInfo(TestElement{"test1"}, 1e8)
	tm.UpdateInfo(TestElement{"test2"}, 2e8)
	tm.UpdateInfo(TestElement{"test3"}, 3e8)
	t.Log(tm.GetAll())
	tm.Delete("test2")
	t.Log(tm.GetAll())
	go tm.UpdateInfo(TestElement{"test4"}, 4e8)
	time.Sleep(2e8)
	t.Log(tm.GetAll())
	go tm.UpdateInfo(TestElement{"test5"}, 5e8)
	go tm.UpdateInfo(TestElement{"test6"}, 6e8)
	go tm.UpdateInfo(TestElement{"test5"}, 5e8)
	go tm.UpdateInfo(TestElement{"test6"}, 6e8)
	go tm.UpdateInfo(TestElement{"test5"}, 5e8)
	go tm.UpdateInfo(TestElement{"test6"}, 6e8)
	go tm.UpdateInfo(TestElement{"test5"}, 5e8)
	go tm.UpdateInfo(TestElement{"test6"}, 6e8)
	time.Sleep(3e8)
	t.Log(tm.GetAll())
	time.Sleep(2.1e8)
	t.Log(tm.GetAll())
	time.Sleep(1e8)
	t.Log(tm.GetAll())

}
