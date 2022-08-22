package gee

// Gee改造ServeHTTP
import (
	"net/http"
)

type HandlerFunc func(c *Context)

// 1. 实例结构体 Engine实现ServeHTTP接口
type Engine struct {
	// 实例的路由表，存储所有路由映射
	router *Router //value是存储的路由，之后通过传入的r ,w 调用
}

// 构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 2. 实现ServeHTTP Func
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 构造Context对象
	c := newContext(w, r)
	// 调用Router的路由器
	engine.router.Handle(c) //将当前会话的请求传入到router进行路由选择
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) GET(pattern string, handlerfunc HandlerFunc) {
	// 绑定方法是：将HandlerFunc添加到engine的路由表中
	engine.router.addRouter("GET", pattern, handlerfunc)
}

func (engine *Engine) POST(pattern string, handlerfunc HandlerFunc) {
	engine.router.addRouter("POST", pattern, handlerfunc)
}
