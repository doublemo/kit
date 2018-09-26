package networks

import "net"

// 定义socket服务器状态
const (
	SocketStateStoped = iota
	SocketStateRunning
)

// SocketConfig 配置
type SocketConfig struct {
	Addr            string                        // 监听地址
	ReadBufferSize  int                           // 读取缓存大小 32767
	WriteBufferSize int                           // 写入缓存大小 32767
	IsGraceful      bool                          // 优雅重启
	CallBack        func(net.Conn, chan struct{}) //
}

// NewSocketConfig 创建socket 配置信息
// 将设定默认值
func NewSocketConfig() *SocketConfig {
	return &SocketConfig{
		ReadBufferSize:  32767,
		WriteBufferSize: 32767,
		IsGraceful:      false,
	}
}
