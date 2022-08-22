package geenote

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //我们使用 roots 来存储每种请求方式的Trie 树根节点
	handlers map[string]HandlerFunc //注册路由
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
// pattern路径切分
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") //会切分出第一个元素是""
	//("/p/hello/:lang", "/")["" "p" "hello" ":lang"]

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 如果有/p/hello/**** 则只存储一个*
				break
			}
		}
	}
	return parts
}

// 路由添加：1. 路由绑定 2.  前缀树添加pattern分支
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method] //get Pots方法是否有根节点
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// pattern查找，返回pattern对应的node：在前缀树中查找(相当于之前在路由表中查找)
// 例如/p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}，/static/css/geektutu.css
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) //等待查找的pattern
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// 在前缀树中搜索，搜索成功则返回对应的node
	n := root.search(searchParts, 0) // get node是0层

	// 找到对应的node
	if n != nil {
		parts := parsePattern(n.pattern) //parts 与searchParts有区别吗？ parts是前缀树中的模糊匹配
		for index, part := range parts {
			// 如果出现模糊匹配， 则将searchParts的精确部分添加到返回值中
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] // 构造：{lang: "go"}，/static/css/geektutu.css
			}
			if part[0] == '*' && len(part) > 1 { //只有*号和不是只有*有什么区别， 只有*匹配所有以*前面开头的pattern
				// /static/css/geektutu.css匹配到/static/*filepath： 把index及其后面的全部都加入进来
				params[part[1:]] = strings.Join(searchParts[index:], "/") //{filepath: "css/geektutu.css"}
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// 根据Context找到对应的处理函数
func (r *router) handle(c *Context) {
	// 获取pattern在前缀树中的结果节点node
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		// 调用对应的处理函数
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
