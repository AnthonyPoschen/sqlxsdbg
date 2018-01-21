package main

/*
	make sure to install this package so it is able to be run from the terminal.
	i.e make sure it is in your path
	go install
	in this directory will put it in the gopath/bin.
*/

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var target *ast.TypeSpec
var info struct {
	packageName string
	StructName  string
	fields      []field
}

type field struct {
	Name string
	Type string
	Tag  string
}

func main() {
	outputFileName := flag.String("o", "", "Set the name of the file that will be saved when generated")
	targetObject := flag.String("t", "", "Set the name of the struct to be parsed")
	flag.Parse()

	if *outputFileName == "" {
		log.Fatal("Please specify ouput file with -o=")
	}
	if *targetObject == "" {
		log.Fatal("Please specify target struct with -t=")
	}
	filename := os.Args[len(os.Args)-1]
	//wddir, _ := os.Getwd()
	//fmt.Println("PWD:", wddir, "- File:", filename, "- Output File:", *outputFileName)
	Ast, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	info.packageName = Ast.Name.Name
	v := visitor{TargetName: *targetObject}
	ast.Walk(v, Ast)
	if target == nil {
		log.Fatal("Unable to find target")
	}
	info.StructName = target.Name.Name
	var w walker
	ast.Walk(w, target)

}

type visitor struct {
	TargetName string
}

type walker int

func (w walker) Visit(n ast.Node) ast.Visitor {
	switch d := n.(type) {
	case *ast.Field:

		spew.Dump(d.Type.(*ast.Ident))
		//typeName := ""
		obj := d.Type.(*ast.Ident)
		fmt.Printf("Name:%v Type:%#v Tag:%v\n", d.Names[0], obj, d.Tag.Value)

		//info.fields = append(info.fields, field{Name: d.Names[0].Name, Type: d.Type, Tag: d.Tag.Value})
	default:
	}
	return w

}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.TypeSpec:
		if d.Name.Name == v.TargetName {
			target = d
			return nil
		}
	default:
	}
	return v
}
