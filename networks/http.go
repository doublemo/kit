package networks

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// WebHandlerFunc 处理函数类型
type WebHandlerFunc func(*WebContext)

// WebHandlersChain 处理函数链
type WebHandlersChain []WebHandlerFunc

// Web HTTP网页服务
type Web struct {
	ctxPool sync.Pool
	router  *Router
}

// allocateContext 创建Ctx
func (web *Web) allocateContext() *WebContext {
	return &WebContext{
		web:      web,
		Request:  NewWebRequest(),
		Response: NewWebResponse(),
	}
}

func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := time.Now()
	ctx := web.ctxPool.Get().(*WebContext)
	ctx.Reset(r, w)
	ctx.Do()
	web.ctxPool.Put(ctx)
	log.Println("[http kit]:", ctx.Response.Status(), ctx.Request.URL.Path, time.Since(s))
}

func (web *Web) Serve(c *WebConfig) error {
	s := &http.Server{
		Addr:           c.Addr,
		Handler:        web,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: c.MaxHeaderBytes,
	}

	return s.ListenAndServe()
}

func (web *Web) ServeTLS(c *WebConfig) error {
	s := &http.Server{
		Addr:           c.Addr,
		Handler:        web,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: c.MaxHeaderBytes,
	}

	return s.ListenAndServeTLS(c.CertFile, c.KeyFile)
}

func (web *Web) ServeUnix(c *WebConfig) error {
	os.Remove(c.Addr)
	listener, err := net.Listen("unix", c.Addr)
	if err != nil {
		return err
	}

	defer listener.Close()
	s := &http.Server{
		Addr:           c.Addr,
		Handler:        web,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: c.MaxHeaderBytes,
	}

	return s.Serve(listener)
}

func (web *Web) Router() *RounterGorup {
	return &RounterGorup{
		basePath: "/",
		root:     true,
		web:      web,
	}
}

func (web *Web) GetRoute(method, path string) *Route {
	return web.router.Find(method, path)
}

func NewWeb() *Web {
	w := &Web{
		router: NewRouter(),
	}

	w.ctxPool.New = func() interface{} {
		return w.allocateContext()
	}
	return w
}
