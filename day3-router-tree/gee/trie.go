package gee

import "strings"

// 【node】
type node struct {
	pattern  string //整个路由， 只有pattern结束的节点才有这个值，其余节点为空
	part     string //当前节点
	children []*node
	isWild   bool //是否是模糊节点
}

// 【insert构建Trie】
func (n *node) insert(pattern string, parts []string, height int) *node {
	// 1. 如果pattern已经到达最后一个，则对应node.pattern设置值
	if len(parts) == height {
		n.pattern = pattern //node能进来就是已经和part匹配上了的，不需要考虑没有匹配上的情况
		return n
	}
	// 2.n.matchChild(part)找到与part匹配的第一个孩子
	part := parts[height]
	child := n.matchChild(part)
	// 3. 没找到对应孩子node则创建一个, 并将这个child节点添加到node的孩子节点中
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	// 4. 继续child.insert()下一个part
	return child.insert(pattern, parts, height+1)
}

// 【search搜索pattern】
func (n *node) search(pattern string, parts []string, height int) *node {
	// 1. 退出:已经到达pattern最后 OR 遇到通配符，都会终止搜索
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 检查最后一个节点是否是终止节点
		if n.pattern != "" {
			return n
		}
		return nil
	}
	// 2. 按层查找
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		// 3. DFS 搜索
		res := child.search(pattern, parts, height+1)

		if res != nil {
			return res
		}
	}
	return nil
}

// 【matchChild匹配part与孩子节点】
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild { //是模糊匹配
			return child
		}
	}
	return nil
}

// 【matchChildren匹配part的所有孩子】
func (n *node) matchChildren(part string) []*node {
	children := []*node{}
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}
