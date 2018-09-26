package networks

import (
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	// LimitRounterWebHandlersChainSize 限制路由处理函数链数据量
	LimitRounterWebHandlersChainSize int = 5
)

type RounterGorup struct {
	basePath string
	root     bool
	web      *Web
	handlers WebHandlersChain
}

// Group 路由组
func (g *RounterGorup) Group(path string, handlers ...WebHandlerFunc) *RounterGorup {
	return &RounterGorup{
		basePath: g.absPath(path),
		handlers: g.mergeWebHandlers(handlers),
		web:      g.web,
	}
}

func (g *RounterGorup) GET(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("GET", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

func (g *RounterGorup) POST(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("POST", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

func (g *RounterGorup) DELETE(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("DELETE", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

func (g *RounterGorup) PATCH(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("PATCH", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

func (g *RounterGorup) PUT(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("PUT", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

func (g *RounterGorup) OPTIONS(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("OPTIONS", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

func (g *RounterGorup) HEAD(urlPath string, handlers ...WebHandlerFunc) {
	assertError(g.web.router.AddRoute("HEAD", g.absPath(urlPath), g.mergeWebHandlers(handlers)))
}

// Any 适配所有HTTP方法
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE
func (g *RounterGorup) Any(urlPath string, handlers ...WebHandlerFunc) {
	abspath := g.absPath(urlPath)
	mergedhandlers := g.mergeWebHandlers(handlers)
	assertError(g.web.router.AddRoute("GET", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("POST", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("PUT", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("PATCH", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("HEAD", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("OPTIONS", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("DELETE", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("CONNECT", abspath, mergedhandlers))
	assertError(g.web.router.AddRoute("TRACE", abspath, mergedhandlers))
}

func (g *RounterGorup) StaticFile(urlPath, filepath string) {
	if strings.Contains(urlPath, ":") || strings.Contains(urlPath, "*") {
		log.Panicln("URL parameters can not be used when serving a static file")
	}

	f := func(c *WebContext) {
		c.File(filepath)
	}

	g.GET(urlPath, f)
	g.HEAD(urlPath, f)
}

func (g *RounterGorup) Static(urlPath, root string) {
	g.StaticFS(urlPath, http.Dir(root))
}

func (g *RounterGorup) StaticFS(urlPath string, fs http.FileSystem) {
	if strings.Contains(urlPath, ":") || strings.Contains(urlPath, "*") {
		log.Panicln("URL parameters can not be used when serving a static folder")
	}

	abspath := g.absPath(urlPath)
	fileServer := http.StripPrefix(abspath, http.FileServer(fs))
	f := func(c *WebContext) {
		fileServer.ServeHTTP(c.Response.GetResponseWriter(), c.Request.GetRequest())
	}

	urlPattern := path.Join(urlPath, "/*filepath")
	g.GET(urlPattern, f)
	g.HEAD(urlPattern, f)
}

func (g *RounterGorup) absPath(urlPath string) string {
	m := path.Join(g.basePath, urlPath)
	if !strings.HasSuffix(m, "/") {
		m += "/"
	}
	return m
}

func (g *RounterGorup) mergeWebHandlers(handlers WebHandlersChain) WebHandlersChain {
	size := len(g.handlers) + len(handlers)
	if size >= LimitRounterWebHandlersChainSize {
		log.Panic("too many handlers")
	}

	mergedHandlers := make(WebHandlersChain, size)
	copy(mergedHandlers, g.handlers)
	copy(mergedHandlers[len(g.handlers):], handlers)
	return mergedHandlers
}
