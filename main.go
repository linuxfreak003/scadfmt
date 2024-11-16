package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ast "github.com/linuxfreak003/scadfmt/ast"
)

func main() {
	var write bool

	flag.BoolVar(&write, "w", false, "Overwrite file")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		panic("invalid args")
	}

	filename := args[0]
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	tree, err := ast.NewAST(file)
	if err != nil {
		log.Fatal(err)
	}
	defer tree.Close()

	out := tree.Format()

	if write {
		err = os.WriteFile(filename, []byte(out), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(out)
	}
}
