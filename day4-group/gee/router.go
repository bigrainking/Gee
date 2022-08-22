package gee

import (
	"net/http"
	"strings"
)

// 【router】
type router struct {
	// 前缀树；路由表
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{map[string]*node{}, map[string]HandlerFunc{}}
}

// 【addRouter】路由添加
func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	// 1. 解析路径parsePattern()
	parts := parsePattern(pattern)
	// 2. roots中插入到前缀树
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{} //根节点为空
	}
	r.roots[method].insert(pattern, parts, 0) //get方法是0层的根节点
	// 3. 添加到路由表
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 【getRouter】路由查询
// 返回查询到的树中的节点， 和模糊查询的值 {lang: /hello}
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	// 1. parsePattern()解析pattern路径以备查询
	searchParts := parsePattern(pattern)

	// 2.r.rootes[method].search() 搜索前缀树中是否有该路径
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(pattern, searchParts, 0)
	// 3. 找到则返回node，如果前缀树中有模糊匹配将具体匹配内容放入param
	if node != nil {
		params := map[string]string{}
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' { //通配符匹配*后面所有部分
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}
		return node, params
	}

	return nil, nil
}

// 【handle】路由选择
func (r *router) handle(c *Context) {
	// 整体流程：1. 在前缀树中查询getRoute()并获取对应的node
	if node, params := r.getRoute(c.Method, c.Path); node == nil {
		c.String(http.StatusNotFound, "404 Note FOUND!")
	} else {
		c.Params = params //模糊匹配的值赋值给会话Context
		//因为 c.Path是请求的具体路由= /p/python   n.patten是树中注册的动态路由
		key := c.Method + "-" + node.pattern
		// 2. 通过handlers路由表调用对应处理函数
		r.handlers[key](c)
	}

}

// 【parsePattern】路径解析
func parsePattern(pattern string) (parts []string) {
	temps := strings.Split(pattern, "/")
	for _, temp := range temps {
		if temp != "" {
			parts = append(parts, temp)
			if temp[0] == '*' { //发现通配符，则只加入一个*
				break
			}
		}
	}
	return
}
