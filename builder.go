package bifurcate

import (
	"errors"
	"fmt"
	"strings"
)

type Builder struct {
	Err  error
	Root Node
}

func (b *Builder) Parse(command string, handler Handler) error {
	commands := strings.Split(strings.TrimSpace(command), " ")
	p := &b.Root
	for _, c := range commands {
		if c[0] == '{' && c[len(c)-1] == '}' {
			item := strings.Split(c[1:len(c)-1], ":")
			if len(item) != 2 {
				return errors.New(fmt.Sprintf("Invalid command %s\n", c))
			}
			literal, valueType := item[0], item[1]
			var t ValueType
			switch valueType {
			case "plain":
				t = Plain
			case "literal":
				t = String
			default:
				return errors.New(fmt.Sprintf("Invalid value type %s\n", valueType))
			}
			node := &Node{
				Literal:   literal,
				Type:      Value,
				ValueType: t,
				Next:      make(map[string]*Node),
			}
			p.AddDefaultNode(node)
			p = node
		} else {
			p = p.AddChildren(&Node{
				Literal: c,
				Type:    Literal,
				Next:    make(map[string]*Node),
			})
		}
	}
	if handler == nil {
		return errors.New("Handler is required\n")
	}
	p.Handler = handler
	return nil
}

func (b *Builder) Build() *Node {
	return &b.Root
}

func NewBuilder() *Builder {
	return &Builder{
		Root: Node{
			Type: Root,
			Next: make(map[string]*Node),
		},
	}
}
