package gee

// 目的：实现一个Handler 处理所有通过端口的请求
import (
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// 1. 实例结构体 Engine实现ServeHTTP接口
type Engine struct {
	// 实例的路由表，存储所有路由映射
	router map[string]HandlerFunc //value是存储的路由，之后通过传入的r ,w 调用
}

// 构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 2. 实现ServeHTTP Func
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 处理进来的request ： 类似switch case，这里借用路由表一个道理
	key := r.Method + "-" + r.URL.Path
	// 如果该请求路径注册过路由，则执行处理函数
	if HandlerFunc, ok := engine.router[key]; ok {
		HandlerFunc(w, r)
	} else {
		fmt.Fprint(w, "404 NOT FOUNT!")
	}
}

// 3. 封装ListenAndServe():
// 让默认走用户自定义的Engine
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// 4. 提供给用户的路由映射处理函数：绑定patter和路由
// 类似HandleFunc：绑定pattern 路由
func (engine *Engine) GET(pattern string, handlerfunc HandlerFunc) {
	// 绑定方法是：将HandlerFunc添加到engine的路由表中
	engine.addRouter("GET", pattern, handlerfunc)
}
func (engine *Engine) POST(pattern string, handlerfunc HandlerFunc) {
	// 绑定方法是：将HandlerFunc添加到engine的路由表中
	engine.addRouter("POST", pattern, handlerfunc)
}

func (engine *Engine) addRouter(method, pattern string, handlerfunc HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handlerfunc
}
