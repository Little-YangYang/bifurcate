package bifurcate

import (
	"errors"
	"fmt"
	"io"
)

type Handler func(ctx *Context, next *Handler) error

type NodeType uint8

const (
	Root NodeType = iota
	Literal
	Value
)

type ValueType uint8

const (
	String ValueType = iota + 1
	Plain
)

type Node struct {
	Next      map[string]*Node
	Type      NodeType
	ValueType ValueType
	Handler   Handler
	Literal   string
	Default   *Node
}

var (
	RouteNotFound = "Route not found\n"
)

// RegisterHandler Register Handler for Node
// @Summary Register Handler for Node. This handler maybe event handler or middleware handler
// @Param handler: the handler for Node
// @Success void
func (n *Node) RegisterHandler(handler Handler) *Node {
	n.Handler = handler
	return n
}

// AddChildren Add Children for Node
// @Summary Add Children for Node
// @Param node: the node to be added
// @Success *Node: the node to be added
func (n *Node) AddChildren(node *Node) *Node {
	n.Next[node.Literal] = node
	return node
}

// AddDefaultNode Add Default for Node
// @Summary Add default child node for Node
// @Param node: the node to be added
// @Success *Node: the node to be added
func (n *Node) AddDefaultNode(node *Node) *Node {
	n.Default = node
	return node
}

func (n *Node) Match(ctx *Context, command *CommandHelper) (c *Context, h Handler, err error) {
	next := ""
	if n.Type == Root {
		if ctx == nil {
			ctx = NewContext()
		}
	}
	if n.Type == Value {
		switch n.ValueType {
		case String:
			val, err := command.GetNext()
			if err != nil && err != io.EOF {
				ctx.Err = err
				return nil, nil, err
			}
			ctx.Data[n.Literal] = val
		case Plain:
			val, err := command.GetNextAll()
			if err != nil && err != io.EOF {
				ctx.Err = err
				return nil, nil, err
			}
			ctx.Data[n.Literal] = val
		}
	}

	next, err = command.GetNext()
	if err != nil && err != io.EOF {
		ctx.Err = err
		return nil, nil, err
	}

	if err == io.EOF {
		if n.Handler != nil {
			return ctx, n.Handler, nil
		} else {
			return nil, nil, errors.New(RouteNotFound)
		}
	}

	// find next node
	if node, ok := n.Next[next]; ok {
		return node.Match(ctx, command)
	} else if n.Default != nil {
		return n.Default.Match(ctx, command.RollBackOne())
	} else if n.Handler != nil {
		return ctx, n.Handler, nil
	} else {
		return nil, nil, errors.New(RouteNotFound)
	}
}

func (n *Node) PrintTree(prefix string, isLast bool) {
	fmt.Print(prefix)
	if isLast {
		fmt.Print("└── ")
		prefix += "    "
	} else {
		fmt.Print("├── ")
		prefix += "│   "
	}

	switch n.Type {
	case Root:
		fmt.Println("Root")
	case Literal:
		fmt.Printf("Literal: %s\n", n.Literal)
	case Value:
		var valueType string
		switch n.ValueType {
		case String:
			valueType = "String"
		case Plain:
			valueType = "Plain"
		}
		fmt.Printf("Value: %s (%s)\n", n.Literal, valueType)
	}

	// Print handler information
	if n.Handler != nil {
		fmt.Printf("%s└── Handler: %p\n", prefix, n.Handler)
	}

	// Print default node
	if n.Default != nil {
		n.Default.PrintTree(prefix, true)
	}

	// Print child nodes
	lastIndex := len(n.Next) - 1
	i := 0
	for _, child := range n.Next {
		child.PrintTree(prefix, i == lastIndex)
		i++
	}
}
