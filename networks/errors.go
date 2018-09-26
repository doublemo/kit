package networks

import "errors"

var (
	ErrorInvalidServerName = errors.New("ErrorInvalidServerName") // 非法的服务器名称
	ErrorServerExist       = errors.New("ErrorServerExist")       // 服务器已经存在
	ErrorServerNonStoped   = errors.New("ErrorServerNonStoped")   // 服务器不是停止状态
)
