package geenote

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)

type (
	// 路由分组
	RouterGroup struct {
		// 该分组的共同前缀 /admin
		prefix string
		// 该分组的中间件：中间件是应用在分组上的
		middlewares []HandlerFunc // support middleware
		// 用于支持分组嵌套： /v1/admin 是 /v1分组的子分组
		parent *RouterGroup // support nesting
		// 便于访问Engine的各种接口：比如Engine.router.addRouter进行路由注册
		// 指针调用engine
		engine *Engine // all groups share a Engine instance
	}

	// Engine实例：？？？为什么需要两个RouterGroup
	Engine struct {
		*RouterGroup
		router router
		//拥有所有RouterGroup的能力
		groups []*RouterGroup // store all groups
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// ===================================路由分组实现===========================

// 创建一个新的路由分组
// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix, //加上parent分组的前缀: /v1  Group(/admin): /v1/admin
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 【分组添加路由】
// 1. 带上共同前缀
// 2. 在路由表中添加完整的路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)

}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
