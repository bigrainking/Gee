package geenote

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
	// 模板相关
	htmlTemplates *template.Template // for html render 文件夹中所有模板文件的集合：将所有的模板加载进内存
	funcMap       template.FuncMap   // for html render 所有的自定义模板渲染函数。
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

// 【中间件注册】
func (group *RouterGroup) Use(h HandlerFunc) {
	// 将中间件添加到group.middlerware中
	group.middlewares = append(group.middlewares, h)
}

// 只有Engine才可以Run  Engine才是实现ServeHTTP的内容
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	c.engine = e //将engine赋值给当前的context
	// 找到r请求所在的group，将group里面的中间件都添加到Context.handlers中
	for _, group := range e.groups {
		// 通过比较RouterGroup与请求的前缀得到
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...) //将所有RouterGroup中的中间件拿出来
		}
	}
	e.router.handle(c)
}

// create static handler 创建文件服务器
// relativePath：用户请求路径
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath) //加上分组的前缀
	// 构建fileServer的handler  http.handler(reqURL, fileServer)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	// 3. 将fileServer转换成HandlerFunc返回， 顺便检查file是否存在
	return func(c *Context) {
		//1) 获取文件：比如请求/static/*filepath param= /assets/js/geektutu.js中的js/geektutu.js
		// ？？？为什么在相对路径获取呢？
		file := c.Param("filepath") // ./static/*filepath对应./static/css/geektutu.css
		// 2） 判断文件是否存在：Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.StatusCode = http.StatusNotFound
			return
		}
		// 启动fileServer的ServeHTTP（http库中已经实现了ServeHTTP：实现了找到文件后如何返回给Writer）
		// 如果不调用ServeHTTP我们需要自己将file传递给writer
		fileServer.ServeHTTP(c.Writer, c.Req) //调用对应的ServeHTTP
	}
}

// serve static files:设置静态文件映射
// @relativePath：用户请求的路径
// @root：文件夹在服务器中的位置
func (group *RouterGroup) Static(relativePath string, root string) {
	// 1. 构建服务器
	// handler := group.createStaticHandler(relativePath, http.Dir(root))
	// 2. 构造实际注册的路径：开发人员注册/static 实际注册的是/static/*filepath
	urlPattern := path.Join(relativePath, "/*filepath")
	// 3. Register GET handlers：注册路由
	// 用户请求 /static/.... 但实际映射到的handler是 /v1/static...
	// group.GET(urlPattern, handler)
	// handler := HandlerFunc(http.StripPrefix( /*absolutePath*/ "/static/", http.FileServer(http.Dir(root /*"./asset"*/)))) //参数2是实际文件夹的路径
	group.GET(urlPattern, handler)
}

// 设置自定义模板函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 加载指定的模板
func (engine *Engine) LoadHTMLGlob(pattern string) {
	// tpl := template.New("").Funcs(engine.funcMap)
	// fmt.Println(tpl.Tree)
	// 解析所有匹配pattern的模板
	// 每个模板对应自己的文件名
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
	// fmt.Println(engine.htmlTemplates)
}
