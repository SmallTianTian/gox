package ast

import (
	"fmt"
	"strings"
)

// FuncAddCall 增加函数调用
// 参数说明：
//
// @father 父函数名称，名称必须是该文件中的一级函数名称，
// 不可是函数内的函数
//
// @underFunc 期望在哪个表达式下面，目前是包含匹配
// 只要文件行包含了该参数的值，即可匹配成功
// 推荐尽可能的多写，进行全匹配
//
// @commond 注释
//
// @fn 函数调用的表达式
//
// @up 如果 @underFunc 参数没有匹配上时，
// 是否期望加在函数的开头，否则将加在末尾
//
// @返回值 是否成功加入函数调用
func (gf *GoFile) FuncAddCall(father, underFunc, commond, fn string, up bool) bool {
	if gf == nil {
		return false
	}

	// 拼装将要加入的字符串
	var addData string
	if commond != "" {
		addData += fmt.Sprintf("// %s\n", commond)
	}
	addData += fn

	// 找到 father 函数申明所在的 Node
	funcNode := gf.lineNode
	fs := fmt.Sprintf("func %s(", father)
	for {
		if funcNode == nil || strings.HasPrefix(funcNode.line, fs) {
			break
		}
		funcNode = funcNode.next
	}

	// 如果没有父节点，则返回 false
	if funcNode == nil {
		return false
	}

	// 没匹配上 underFunc 时使用的节点
	var rescueNode *node

	// 向下迭代 father 函数申明所在的 Node
	// 进行 underFunc 匹配，以求能加入到其中
	stop := funcNode.IterDown(func(n *node) bool {
		rescueNode = n
		if !strings.Contains(n.line, underFunc) {
			return false
		}

		n.Add(&node{line: addData})
		return true
	})

	// 没有主动暂停，代表没有匹配上 underFunc
	if !stop {
		// 如果是要加入到函数的最前端，
		// 则将 node 重置为 father 函数申明所在的 Node
		if up {
			rescueNode = funcNode
		}
		rescueNode.Add(&node{line: addData})
	}
	return true
}
