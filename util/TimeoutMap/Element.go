package TimeoutMap

type Element interface {
	GetID() string
	NewAddedHandler()
	TimeoutHandler()
}
