package bifurcate

import (
	"io"
	"strings"
)

type CommandHelper struct {
	prev    int
	next    int
	Command []rune
}

func (c *CommandHelper) GetCommand() string {
	return string(c.Command)
}

func (c *CommandHelper) RollBackOne() *CommandHelper {
	i := min(c.prev, len(c.Command)-1)
	for ; i >= 0; i-- {
		if c.Command[i] == ' ' {
			c.prev = i
			break
		}
	}
	if i < 0 {
		c.prev = -1
	}
	return c
}

func (c *CommandHelper) GetNext() (string, error) {

	defer func() {
		c.prev = c.next
	}()

	// judge if command have continuous space flag
	literalFlag := false
	i := c.prev + 1
	for ; i < len(c.Command); i++ {
		if c.Command[i] != ' ' {
			literalFlag = true
		}
		if c.Command[i] == ' ' && literalFlag {
			c.next = i
			break
		}
	}
	if i <= c.prev+1 {
		return "", io.EOF
	}
	if i == len(c.Command) {
		c.next = i
		return string(c.Command[c.prev+1:]), nil
	}
	return string(c.Command[c.prev+1 : c.next]), nil
}

func (c *CommandHelper) GetNextAll() (string, error) {
	if c.prev == len(c.Command) {
		return "", io.EOF
	}
	return string(c.Command[c.prev+1:]), nil
}

func NewCommandHelper(command string) *CommandHelper {
	command = strings.Trim(command, " ")
	return &CommandHelper{
		prev:    -1,
		next:    0,
		Command: []rune(command),
	}
}
