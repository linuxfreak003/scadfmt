package ast

import (
	"context"
	"fmt"
	"strings"
	"time"

	ts_openscad "github.com/linuxfreak003/tree-sitter-openscad/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

type AST struct {
	tree            *sitter.Tree
	src             []byte
	sb              strings.Builder
	indentLevel     int
	tempIndentLevel []int
	indent          string
}

func NewAST(src []byte) (*AST, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(sitter.NewLanguage(ts_openscad.Language()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tree := parser.ParseCtx(ctx, src, nil)
	root := tree.RootNode()

	if root.HasError() {
		return nil, fmt.Errorf("Error in File")
	}

	return &AST{
		tree:            tree,
		src:             src,
		sb:              strings.Builder{},
		indent:          "  ",
		tempIndentLevel: []int{0},
	}, nil
}

func (a *AST) SetIndent(s string) {
	a.indent = s
}

func (a *AST) Indent() {
	a.indentLevel++
	a.tempIndentLevel = append(a.tempIndentLevel, 0)
}

func (a *AST) Dedent() {
	a.indentLevel--
	if a.indentLevel < 0 {
		a.indentLevel = 0
	}
	a.tempIndentLevel = a.tempIndentLevel[:len(a.tempIndentLevel)-1]
}

func (a *AST) IndentLevel() int {
	return a.indentLevel + sum(a.tempIndentLevel)
}

func sum(a []int) (s int) {
	for _, v := range a {
		s += v
	}
	return
}

func (a *AST) TempIndent() {
	a.tempIndentLevel[a.indentLevel]++
}
func (a *AST) TempReset() {
	a.tempIndentLevel[a.indentLevel] = 0
}

func (a *AST) Format() string {
	root := a.tree.RootNode()

	var visit func(n *sitter.Node, name string, depth int)
	visit = func(n *sitter.Node, name string, depth int) {
		//spaces := strings.Repeat("|", depth)
		//fmt.Printf("%sName:%s Kind:%s\n", spaces, n.GrammarName(), n.Kind())
		//fmt.Printf("%s%s (%s)", spaces, n.Kind(), name)
		if n.ChildCount() == 0 {
			text := n.Utf8Text(a.src)
			a.sb.WriteString(text)
			//fmt.Printf(" - [%s]", text)
		}
		//fmt.Printf("\n")
		for i := 0; i < int(n.ChildCount()); i++ {
			child := n.Child(uint(i))
			a.Before(child)
			visit(child, n.FieldNameForChild(uint32(i)), depth+1)
			a.After(child)

		}
	}

	visit(root, "root", 0)

	return a.sb.String()
}

func (a *AST) Before(node *sitter.Node) {
	switch node.Kind() {
	case "transform_chain":
		if node.Parent().Kind() == "transform_chain" {
			a.TempIndent()
			a.sb.WriteString("\n")
		}
		a.sb.WriteString(strings.Repeat(a.indent, a.IndentLevel()))
	case "module", "for_block":
		a.sb.WriteString(strings.Repeat(a.indent, a.IndentLevel()))
	case "{":
		a.sb.WriteString(" ")
	case "}":
		a.Dedent()
		a.sb.WriteString(strings.Repeat(a.indent, a.IndentLevel()))
	default:
	}
}

func (a *AST) After(node *sitter.Node) {
	switch node.Kind() {
	case "use_statement":
		a.sb.WriteString("\n")
		if node.NextSibling().Kind() != "use_statement" {
			a.sb.WriteString("\n")
		}
	case ",":
		a.sb.WriteString(" ")
	case "module":
		a.sb.WriteString(" ")
	case "{":
		a.Indent()
		a.sb.WriteString("\n")
	case ";":
		a.TempReset()
		a.sb.WriteString("\n")
	case "}":
		a.TempReset()
		a.sb.WriteString("\n")
		if g := node.Parent().Parent(); g != nil && g.NextSibling() != nil && g.NextSibling().Kind() != "}" {
			a.sb.WriteString("\n")
		}
	}
}

func (a *AST) Close() {
	if a.tree != nil {
		a.tree.Close()
	}
}
