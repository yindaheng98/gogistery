package base

type SenderInfo interface {
	Send(addr string) (ReceiverInfo, error)
}
