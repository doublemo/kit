package networks

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
)

type WebContext struct {
	web      *Web
	Request  *WebRequest
	Response WebResponseWriter
	abort    bool // 结束当前运行
}

func (wctx *WebContext) Method() string {
	return wctx.Request.Method
}

func (wctx *WebContext) GetHeader(key string, val ...string) string {
	data := wctx.Request.Header.Get(key)
	if len(data) < 1 && len(val) > 0 {
		return val[0]
	}

	return data
}

func (wctx *WebContext) SetHeader(k, v string) {
	wctx.Response.Header().Set(k, v)
}

func (wctx *WebContext) Do() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Runtime Error:", err)
			log.Printf("%s\n", debug.Stack())
			wctx.RenderError(http.StatusInternalServerError, fmt.Errorf("Runtime Error: %v ", err))
		}
	}()

	route := wctx.web.GetRoute(wctx.Method(), wctx.Request.URL.Path)
	if route == nil {
		wctx.RenderError(http.StatusNotFound, fmt.Errorf("nomatch: %s %s", wctx.Method(), wctx.Request.URL.Path))
		return
	}

	for _, h := range route.handlers {
		if h == nil {
			continue
		}

		h(wctx)
		if wctx.IsAborted() {
			break
		}
	}
}

func (wctx *WebContext) RenderError(status int, err error) {
	if status == 0 || status == http.StatusOK {
		status = http.StatusInternalServerError
	}

	wctx.Abort()
	wctx.Response.WriteStatus(status)
	wctx.Response.WriteContentType("")
	wctx.Response.WriteString(err.Error())
}

func (wctx *WebContext) RenderText(msg string) {
	wctx.Response.WriteContentType("text/plain; charset=utf-8")
	wctx.Response.WriteString(msg)
}

func (wctx *WebContext) Abort() {
	wctx.abort = true
}

func (wctx *WebContext) IsAborted() bool {
	return wctx.abort
}

func (wctx *WebContext) Reset(r *http.Request, w http.ResponseWriter) {
	wctx.Request.Apply(r)
	wctx.Response.Apply(w)
	wctx.abort = false
}

// File 文件服务
func (wctx *WebContext) File(filepath string) {
	http.ServeFile(wctx.Response.GetResponseWriter(), wctx.Request.GetRequest(), filepath)
}

// IsWebsocket 检查是否为websocket
func (wctx *WebContext) IsWebsocket() bool {
	connection := strings.ToLower(wctx.GetHeader("Connection"))
	upgrade := strings.ToLower(wctx.GetHeader("Upgrade"))

	if strings.Contains(connection, "upgrade") && upgrade == "websocket" {
		return true
	}

	return false
}

// ClientIP 获取客户端IP信息
func (wctx *WebContext) ClientIP() string {
	clientIP := wctx.Request.Header.Get("X-Forwarded-For")
	if idx := strings.IndexByte(clientIP, ','); idx >= 0 {
		clientIP = clientIP[0:idx]
	}

	clientIP = strings.TrimSpace(clientIP)
	if clientIP != "" {
		return clientIP
	}

	clientIP = strings.TrimSpace(wctx.Request.Header.Get("X-Real-Ip"))
	if clientIP != "" {
		return clientIP
	}

	remoteAddr := strings.TrimSpace(wctx.Request.RemoteAddr)
	if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return ip
	}

	return ""
}

func (wctx *WebContext) ContentType() string {
	return wctx.Request.ContentType
}
