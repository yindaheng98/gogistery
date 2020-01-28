package Single

type callbacks struct {
	Started func()
	Stopped func()
}

func newCallbacks() *callbacks {
	return &callbacks{func() {}, func() {}}
}
