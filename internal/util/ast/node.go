package ast

import "strings"

type node struct {
	line string
	// 前、后一个节点
	pre, next *node
	// 上、下一个区间
	up, down []*node
}

func newNode(lines []string) *node {
	if len(lines) == 0 {
		return nil
	}

	first := &node{line: lines[0]}
	var waitMatchNode []*node
	cur := first
	for _, line := range lines[1:] {
		n := &node{line: line, pre: cur}
		cur.next = n
		cur = n

		lk := strings.Count(line, "{")
		rk := strings.Count(line, "}")
		if lk == rk {
			continue
		}

		diff := lk - rk
		if diff > 0 {
			for i := 0; i < diff; i++ {
				waitMatchNode = append(waitMatchNode, n)
			}
			continue
		}
		for i := diff; i < 0; i++ {
			l := waitMatchNode[len(waitMatchNode)-1]
			waitMatchNode = waitMatchNode[:len(waitMatchNode)-1]

			l.down = []*node{n}
			n.up = append(n.up, l)
		}
	}
	return first
}

// HasNext 是否还有子节点
func (n *node) HasNext() bool {
	return n.next != nil
}

// HasPre 是否还有父节点
func (n *node) HasPre() bool {
	return n.next != nil
}

// Add 新增加节点，多次调用，类似入栈。
//
// 例如： Add(1).Add(2).Add(3)
//
// 效果：3-2-1
func (n *node) Add(new *node) *node {
	if new == nil {
		return n
	}

	nn := n.next

	new.next = nn
	new.pre = n
	n.next = new

	if nn != nil {
		nn.pre = new
	}

	return n
}

// TillDown 一直向下取值，直到下区间
//
// 返回值：调用的函数是否主动暂停
func (n *node) IterDown(f func(*node) bool) bool {
	d := n.down[len(n.down)-1]
	cur := n.next
	for {
		if cur == nil || cur == d {
			return false
		}
		if f(cur) {
			return true
		}
		cur = cur.next
	}
}

// TillUp 一直向上取值，直到上区间
//
// 返回值：调用的函数是否主动暂停
func (n *node) IterUp(f func(*node) bool) bool {
	p := n.up[len(n.up)-1]
	cur := n.pre
	for {
		if cur == nil || cur == p {
			return false
		}
		if f(cur) {
			return true
		}
		cur = cur.pre
	}
}

func (n *node) String() string {
	if n == nil {
		return ""
	}

	var sb strings.Builder
	cur := n
	for cur != nil {
		sb.WriteString(cur.line)
		sb.WriteString("\n")
		cur = cur.next
	}
	return sb.String()
}
