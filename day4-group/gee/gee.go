package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)

// 【分组路由RouterGroup】
type RouterGroup struct {
	// 共同前缀
	prefix string //不导出，外面用不到
	// 中间件
	middlewares []HandlerFunc
	// parent
	parent *RouterGroup
	// Engine的调用:包含Engine的所有接口，Router相当于Engine下面的子分组
	engine *Engine
}

// 【Engine对象】
type Engine struct {
	router *router
	// RouterGroup属性： 父分组拥有子分组的全部属性功能
	*RouterGroup
	// 包含所有子分组的功能
	groups []*RouterGroup
}

// 【Engine构造函数:全局构造函数，创建一个实例】
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = append(engine.groups, engine.RouterGroup)
	return engine
}

// 【创建路由分组Group】提供给开发人员使用创建一个分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine //group是可以获取到全局的Engine的， 将全局的Engine赋值给该新建的路由分组
	routerGroup := &RouterGroup{
		prefix: group.prefix + prefix, //参数传入的只是子分组下面的前缀，加上parent分组的前缀: /v1  Group(/admin): /v1/admin
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, routerGroup)
	return routerGroup
}

// 【路由分组：添加路由条目(addRouter)】
func (group *RouterGroup) addRouter(method, pattern string, handler HandlerFunc) {
	// 1. 完善路径：此处只是给出的路由分组下的路径，需要加上父分组的前缀
	patternall := group.prefix + pattern
	log.Printf("Route %4s - %s", method, patternall) // 路由分组下注册 log记录
	// 2. 调用router的路由注册函数
	group.engine.router.addRouter(method, patternall, handler)
}

// 【路由分组：Get、POST】
// 分组下面传入的所有pattern都是建立在共同前缀下的
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

// 只有Engine才可以Run  Engine才是实现ServeHTTP的内容
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	// 插入点
	// 允许用户使用自己定义的中间件做一些额外的处理，例如记录日志等，以及对Context进行二次加工。
	e.router.handle(c)
	// 另外通过调用(*Context).Next()函数，
	// 中间件可等待用户自己定义的 Handler处理结束后，再做一些额外的操作，例如计算本次处理所用时间等。
	// 即 Gee 的中间件支持用户在请求被处理的前后，做一些额外的操作。
}
