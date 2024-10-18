package main

import (
	ast "github.com/linuxfreak003/scadfmt/ast"
)

func main() {
	tree := &ast.AST{}
	tree.Parse([]byte("const foo = 1 + 2"))
}
