package Sender

type StatusType uint32

const (
	STATUS_Disconnected StatusType = 0
	STATUS_Retrying     StatusType = 1
	STATUS_Connected    StatusType = 2
)
