package ast

import (
	"fmt"
	"strings"
)

// StructAddField 结构体里面增加字段
// 参数说明：
//
// @fathers 父结构体名称列表，必须从文件的一级父结构体开始，
// 没有找到的结构体，将自动创建
//
// 例如: []string{"Config", "SubConfig"} 将匹配到
//
// type Config struct {
//		SubConfig struct { // 匹配到这个子结构体
//		}
// }
//
// @commond 注释
//
// @fn 添加字段内容，例如:
// 1. 添加年龄字段 "Age int"
// 2. 添加子结构体 "ChildConfig struct"
//
// 只创建 struct，更建议加入到 fathers 字段中，
// 让本字段为空，例如：[]string{"Config", "SubConfig", "ChildConfig"}
//
// @isStruct 添加的字段是否是结构体
//
// @返回值 是否成功加入字段
func (gf *GoFile) StructAddField(fathers []string, commond, fn string) bool { // nolint
	if gf == nil {
		return false
	}
	if len(fathers) == 0 {
		return false
	}

	curNode := gf.lineNode
	structName := fathers[0]

	// 获取一级父结构体
	for {
		// 如果到末尾还没有匹配上，则返回失败
		if curNode == nil {
			return false
		}
		if strings.HasSuffix(curNode.line, fmt.Sprintf("type %s struct {", structName)) {
			break
		}
		curNode = curNode.next
	}

	// 无限逼近最内部的父节点
	for {
		fathers = fathers[1:]
		if len(fathers) == 0 {
			break
		}
		initiative := curNode.IterDown(func(n *node) bool {
			if strings.HasSuffix(n.line, fmt.Sprintf("%s struct {", fathers[0])) {
				curNode = n
				return true
			}
			return false
		})
		// 如果不主动退出，代表没有匹配上，
		// 则退出逼近模式
		if !initiative {
			break
		}
	}

	// 如果剩下的还有父节点，则依次创建
	for _, fat := range fathers {
		u := &node{line: fmt.Sprintf("%s struct {", fat)}
		d := &node{line: "}", up: []*node{u}}
		u.down = append(u.down, d)
		curNode.Add(d).Add(u)

		curNode = u
	}

	// 如果字段为空，则直接返回
	// 当只创建 struct 时，存在这种情况
	if fn == "" {
		return true
	}

	// 这里才是真正创建字段的地方
	fn = strings.TrimSpace(fn)
	isStruct := strings.HasSuffix(fn, " struct")

	// 如果不是 struct，则简单的创建
	if !isStruct {
		content := fn
		if commond != "" {
			content += " // " + commond
		}
		curNode.Add(&node{line: content})
		return true
	}

	// 如果是 struct
	var content string
	if commond != "" {
		content = fmt.Sprintf("// %s\n", commond)
	}
	content += fn + "{"

	u := &node{line: content}
	d := &node{line: "}", up: []*node{u}}
	u.down = append(u.down, d)
	curNode.Add(d).Add(u)
	return true
}
