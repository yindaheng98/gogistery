package Error

type ConnectionErrorCode uint32

const (
	CONNECTION_ERROR_InitFailed          ConnectionErrorCode = 0 //初始化阶段出错
	CONNECTION_ERROR_RetryFailed         ConnectionErrorCode = 1 //重试连接时出错
	CONNECTION_ERROR_ConnectionInterrupt ConnectionErrorCode = 2 //本来连的好好的突然断了
)

type ConnectionError struct {
	error
	Code ConnectionErrorCode
}

func NewConnectionError(err error, code ConnectionErrorCode) ConnectionError {
	return ConnectionError{err, code}
}
