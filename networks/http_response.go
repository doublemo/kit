package networks

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

type WebResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// 构造输出
	Apply(http.ResponseWriter)

	// 返回http当前状态码
	Status() int

	// 返回当前输出内容长度
	Size() int

	// 写入一个字第串到输出内容中
	WriteString(string) (int, error)

	// 写入当前http状态码
	WriteStatus(int)

	// 写入当前内容mime
	WriteContentType(string)

	// 返回当前是否准备就绪
	ReadyedWritten() bool

	// 获取原有Writer
	GetResponseWriter() http.ResponseWriter
}

type WebResponse struct {
	http.ResponseWriter
	status      int
	size        int
	contentType string
}

func (w *WebResponse) GetResponseWriter() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *WebResponse) Apply(rw http.ResponseWriter) {
	w.ResponseWriter = rw
	w.Reset()
}

func (w *WebResponse) Write(data []byte) (int, error) {
	w.writeHeader()
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *WebResponse) WriteString(s string) (int, error) {
	w.writeHeader()
	n, err := io.WriteString(w.ResponseWriter, s)
	w.size += n
	return n, err
}

func (w *WebResponse) WriteStatus(code int) {
	if code > 0 && w.status != code {
		w.status = code
	}
}

func (w *WebResponse) WriteContentType(s string) {
	if len(s) > 0 && w.contentType != s {
		w.contentType = s
	}
}

func (w *WebResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}

	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *WebResponse) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *WebResponse) Flush() {
	w.writeHeader()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *WebResponse) ReadyedWritten() bool {
	return w.size != -1
}

func (w *WebResponse) Status() int {
	return w.status
}

func (w *WebResponse) Size() int {
	return w.size
}

func (w *WebResponse) writeHeader() {
	if !w.ReadyedWritten() {
		w.size = 0
		w.ResponseWriter.Header().Set("Content-Type", w.contentType)
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *WebResponse) Reset() {
	w.status = http.StatusOK
	w.size = -1
	w.contentType = ""
}

func NewWebResponse() *WebResponse {
	return &WebResponse{
		status: http.StatusOK,
		size:   -1,
	}
}
