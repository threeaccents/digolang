package compiler

import (
	"github.com/threeaccents/digolang/ast"
	"github.com/threeaccents/digolang/code"
	"github.com/threeaccents/digolang/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node *ast.Node) error {
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: code.Instructions{},
		Constants:    []object.Object{},
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
