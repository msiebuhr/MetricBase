package graphiteParser

import (
	"fmt"
	"strings"
)

type NodeType uint8

const (
	NODE_UNKNOWN  NodeType = iota
	NODE_ROOT              // Internal to parser
	NODE_FUNCTION          // foo(...)
	NODE_METRIC            // statsd.appname.*.upper_90
	NODE_NUMBER            // 12.2 / 10
	NODE_STRING            // "Bytes per second"
)

// Node
type Node struct {
	Name   string
	Char   int
	Type   NodeType
	Args   []*Node
	parent *Node
}

func NewNode(name string, char int, ntype NodeType) *Node {
	return &Node{
		Name: name,
		Char: char,
		Type: ntype,
		Args: make([]*Node, 0),
	}
}

func (p *Node) addArgument(arg *Node) {
	p.Args = append(p.Args, arg)
}

func (p *Node) String() string {
	// Print functions in a magic way
	if p.Type == NODE_FUNCTION {
		argStrings := make([]string, len(p.Args))
		for i := range p.Args {
			argStrings[i] = p.Args[i].String()
		}
		return fmt.Sprintf("%s(%s)", p.Name, strings.Join(argStrings, ", "))
	}

	// Or just the name
	return p.Name
}
