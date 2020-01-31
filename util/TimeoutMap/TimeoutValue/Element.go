package TimeoutValue

type Element interface {
	NewAddedHandler()
	TimeoutHandler()
}
