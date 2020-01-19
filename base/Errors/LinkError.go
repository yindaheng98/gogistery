package Errors

import "gogistery/base"

//用于指示在“从哪个sender到哪个receiver的发送过程中发生了什么error”
type LinkError struct {
	error
	link base.LinkInfo
}

func NewLinkError(err error, linkInfo base.LinkInfo) LinkError {
	return LinkError{error: err, link: linkInfo}
}

func (e *LinkError) Pair() base.LinkInfo {
	return e.link
}
