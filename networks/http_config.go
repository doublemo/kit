package networks

import (
	"time"
)

// WebConfig Web服务器配置
type WebConfig struct {
	// TCP address to listen on, ":http" if empty
	Addr string

	// 读取内容超时
	//
	// 如果浏览器请求时还好,可以保护客户端能及时关闭连接
	// 当客户端使用某些库时不能保证关闭连接时，如果造成服务器连接越来越多不能及时释放
	ReadTimeout time.Duration

	// 写入内容超时
	WriteTimeout time.Duration

	// 限制解析HTTP头大小
	MaxHeaderBytes int

	// SSL证书
	CertFile string

	// SSL证书
	KeyFile string
}

// NewWebConfig 创建配置,并指定默认值
func NewWebConfig(addr string) *WebConfig {
	return &WebConfig{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
