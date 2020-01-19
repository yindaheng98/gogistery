package base

type LinkInfo struct {
	sender   SenderInfo
	receiver ReceiverInfo
}

func NewLinkInfo(senderInfo SenderInfo, receiverInfo ReceiverInfo) LinkInfo {
	return LinkInfo{sender: senderInfo, receiver: receiverInfo}
}

func (p LinkInfo) SenderInfo() SenderInfo {
	return p.sender
}

func (p LinkInfo) ReceiverInfo() ReceiverInfo {
	return p.receiver
}
