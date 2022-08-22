package gee

import "log"

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
func NEW() *Engine {
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
	log.Printf("Route %4s - %s", method, pattern) // 路由分组下注册 log记录
	// 2. 调用router的路由注册函数
	group.engine.router.addRouter(method, patternall, handler)
}

// 【路由分组：Get、POST】
// 分组下面传入的所有pattern都是建立在共同前缀下的
func (group *RouterGroup) GET(pattern, handler HandlerFunc) {

}

// 只有Engine才可以Run  Engine才是实现ServeHTTP的内容
