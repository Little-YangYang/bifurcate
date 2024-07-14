package bifurcate

import (
	"fmt"
	"log"
	"testing"
)

func TestBuilder(t *testing.T) {
	builder := NewBuilder()
	handler := func(ctx *Context, next *Handler) error {
		fmt.Println(ctx.Data)
		return nil
	}

	// Test case 1: Parse a valid command
	err := builder.Parse("foo bar {id:literal}", handler)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test case 2: Parse an invalid command
	err = builder.Parse("foo bar {id}", handler)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	// Test case 3: Parse without a handler
	err = builder.Parse("foo bar {id:literal}", nil)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestNode_Match(t *testing.T) {
	builder := NewBuilder()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handler := func(ctx *Context, next *Handler) error {
		return nil
	}
	err := builder.Parse("foo bar {id:literal}", handler)
	if err != nil {
		t.Log(err)
	}
	root := builder.Build()
	//fmt.Println(root)
	root.PrintTree("", true)
	// Test case 1: Match a valid command
	ctx, h, err := root.Match(nil, NewCommandHelper("foo bar 123"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ctx.GetData("id") != "123" {
		t.Errorf("Expected id to be '123', but got '%v'", ctx.GetData("id"))
	}
	if h == nil {
		t.Error("Expected a handler, but got nil")
	}

	// Test case 2: Match an invalid command
	ctx, h, err = root.Match(nil, NewCommandHelper("foo baz"))
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
