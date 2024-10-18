package ast

import (
	"fmt"

	ts_openscad "github.com/linuxfreak003/tree-sitter-openscad/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type AST struct{}

func (a *AST) Parse(src []byte) error {
	code := []byte("cube(10);")

	parser := tree_sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(tree_sitter.NewLanguage(ts_openscad.Language()))

	tree := parser.Parse(code, nil)
	defer tree.Close()

	root := tree.RootNode()
	fmt.Println(root.ToSexp())

	return nil
}
