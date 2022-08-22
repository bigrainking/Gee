package gee

// 将router单独拎出来
import (
	"net/http"
)

// 路由器
type Router struct {
	handlers map[string]HandlerFunc
}

// 构造函数：创建一个Router
func newRouter() *Router { //供内部使用需要小写
	return &Router{handlers: make(map[string]HandlerFunc)}
}

func (router *Router) addRouter(method, pattern string, handlerfunc HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handlerfunc
}

// 提供给ServeHTTP的路由选择函数
func (router *Router) Handle(c *Context) {
	// 处理进来的request ： 类似switch case，这里借用路由表一个道理
	key := c.Method + "-" + c.Path
	// 如果该请求路径注册过路由，则执行处理函数
	if HandlerFunc, ok := router.handlers[key]; ok {
		HandlerFunc(c)
	} else {
		// 构造返回内容
		c.String(http.StatusNotFound, "404 NOT FOUND！")
	}
}
