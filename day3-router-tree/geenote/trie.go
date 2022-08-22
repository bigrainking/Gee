package geenote

import "strings"

// 节点，除了最后一个节点以外其他节点都不会有pattern
// 比如/p/:lang/doc是注册了路由的，同时他会每部分拆分成一个node
// 当查询路径 /p/python 时候这个路径并没有注册路由，但是我们可以匹配到 /p/:lang，如何告诉它你没有匹配到路径呢？
// /:lang节点不是结尾节点pattern值为空，
// 对于/p/python已经匹配完毕但是匹配到的最后一个/:lang节点pattern值为空， 因此匹配失败
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，同一层的所有子节点例如 [doc, tutorial, intro](处在同一层)
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 返回node的children中第一个匹配part成功的节点：用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild { //模糊匹配一定会匹配成功
			return child
		}
	}
	return nil
}

// 返回node孩子节点中所有匹配part成功的节点，用于查找
// 比如 /hello  匹配成功 /hello /:name /*
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 路径插入: 在node的孩子层插入节点
// 【调用】r.roots[method].insert(pattern, parts, 0)：在0层(Get节点)后插入pattern路径
// height：当前节点在树中的层数
// 递归查找每一层节点，如果没有找到则插入一个child节点
// 注意：如果原树中有路径 /p/:lang/dox
// 现在要插入 /p/python, 在查找过程中发现匹配到了，但是/:langpattern值为空 ： 则要在第二层插入/python节点
// 并且python节点的pattern值 = "/p/python"
func (n *node) insert(pattern string, parts []string, height int) {
	// 1. 如果当前将需要匹配的pattern匹配完了，则匹配结束，并将当前node.pattern 设置：parts[0]匹配的是height=1层的node
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height] //获取pattern的第height个值
	// 2. 对于本node的孩子中是否有与part的匹配的：part匹配node
	child := n.matchChild(part)
	// 3. 如果没有找到匹配的，则创建一个node
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 添加到孩子集合中
		n.children = append(n.children, child)
	}
	// 4. 继续匹配node.child和下一个part
	child.insert(pattern, parts, height+1)
}

// 路径搜索：返回匹配到的最后一个节点(带有pattern值的节点)
// 每层搜索，直到搜索到有pattern的节点并匹配得上。 如果到了最后一层or匹配到*
//退出规则是，匹配到了*，匹配失败，????
func (n *node) search(parts []string, height int) *node {
	// 1. 退出： 已经到达pattern最后一个part OR 遇到通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" { //如果该node不是终结节点则匹配失败： /p/:lang/index 匹配 /p/python 是失败的
			return nil
		}
		return n
	}
	// 2. 匹配node节点的孩子与part
	part := parts[height]
	children := n.matchChildren(part)
	// 3. 对于匹配成功的每个孩子节点，分别进入对应支路去搜索
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
