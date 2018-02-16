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
	"reflect"
	"strings"
	"text/template"
)

// Internal Globals

var target *ast.TypeSpec
var info struct {
	PackageName  string
	StructName   string
	DatabaseName string
	TableName    string
	Fields       []field
}

// Types

type field struct {
	Name      string
	Type      string
	Tag       string
	IsPointer bool
}

type targetFinder struct {
	TargetName string
}

func (v targetFinder) Visit(n ast.Node) ast.Visitor {
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

type targetStructWalker int

func (w targetStructWalker) Visit(n ast.Node) ast.Visitor {
	switch d := n.(type) {
	case *ast.Field:
		switch b := d.Type.(type) {
		case *ast.Ident:
			obj := b
			info.Fields = append(info.Fields, field{Name: d.Names[0].Name, Type: obj.Name, Tag: d.Tag.Value})
		case *ast.StarExpr:
			obj := b.X.(*ast.Ident)
			info.Fields = append(info.Fields, field{Name: d.Names[0].Name, Type: obj.Name, Tag: d.Tag.Value, IsPointer: true})
		default:
		}
	default:
	}
	return w

}

// Main entry
func main() {
	targetObject := flag.String("t", "", "Set the name of the struct to be parsed")
	dbName := flag.String("db", "", "Database Name for sqlx")
	tbName := flag.String("tb", "", "Table Name within the database to associate target struct with")
	help := flag.Bool("h", false, "Shows this Help Dialogue")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *targetObject == "" {
		log.Fatal("Target struct missing '-t=VALUE', use -h for help")
	}

	if *dbName == "" {
		log.Fatal("Database Name missing '-db=VALUE', use -h for help")
	}

	if *tbName == "" {
		log.Fatal("Table Name missing '-tb=VALUE', use -h for help")
	}
	info.DatabaseName = *dbName
	info.TableName = *tbName
	filename := os.Args[len(os.Args)-1]
	outputfilename := filename[0:len(filename)-3] + "_gen.go"
	Ast, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Fatal("Unable to parse file:", err)
	}

	info.PackageName = Ast.Name.Name
	tf := targetFinder{TargetName: *targetObject}
	ast.Walk(tf, Ast)
	if target == nil {
		log.Fatal("Unable to find target in file")
	}
	info.StructName = target.Name.Name
	var tsw targetStructWalker
	ast.Walk(tsw, target)

	funcMap := template.FuncMap{
		"constants":     buildConstants,
		"lowerCase":     lowerCaseFirst,
		"getFunc":       getFunc,
		"getMultiFunc":  getMultiFunc,
		"saveFunc":      saveFunc,
		"saveMultiFunc": saveMultiFunc,
		"newFunc":       newFunc,
	}

	tmpl, err := template.New("test").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatal(err)
	}
	// delete the file if it exists, we don't give a shit about an error here
	_ = os.Remove(outputfilename)
	//open file and pass that instead of stdout
	file, err := os.OpenFile(outputfilename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Failed to open final file:%s", err)
	}
	err = tmpl.Execute(file, info)
	if err != nil {
		log.Fatalf("Template Failed to execute: %s", err)
	}

}

func lowerCaseFirst(s string) string {
	r := strings.ToLower(string(s[0]))
	r += s[1:]
	return r
}

func buildConstants() (result string) {
	structName := lowerCaseFirst(info.StructName)
	result += "\t" + info.StructName + "SearchTypeLIKE " + structName + "SearchType = \"LIKE\"\n"
	result += "\t" + info.StructName + "SearchTypeEQUAL " + structName + "SearchType = \"=\"\n"
	result += "\t" + structName + "TableName string = \"" + info.TableName + "\"\n\n"
	for k, v := range info.Fields {
		cleanRawTags := v.Tag[1 : len(v.Tag)-1]
		tag := reflect.StructTag(cleanRawTags)

		dbTag, ok := tag.Lookup("db")
		if ok != true {
			fmt.Println("DB Tag not found", tag.Get("db"))
			continue
		}
		if k != 0 {
			result += "\n"
		}
		result += "\t" + info.StructName + "Field" + v.Name + " " + structName + "Field = \"" + dbTag + "\""
	}
	return result
}

func getFunc() (result string) {
	fieldStruct := lowerCaseFirst(info.StructName) + "Field"
	tableName := lowerCaseFirst(info.StructName) + "TableName"
	result += fmt.Sprintf("func %sGet( db *sqlx.DB, key %s, value string) (%s, error) {\n", info.StructName, fieldStruct, info.StructName)
	result += "	var result " + info.StructName + "\n"
	result += `	statement := fmt.Sprintf("SELECT * from %s.%s where %s=?", "` + info.DatabaseName + `", ` + tableName + `, key)
	return result, db.Unsafe().Get(&result,statement,value)
}`
	return

}

func getMultiFunc() (result string) {
	fieldStruct := lowerCaseFirst(info.StructName) + "Field"
	searchSruct := lowerCaseFirst(info.StructName) + "SearchType"
	tableName := lowerCaseFirst(info.StructName) + "TableName"
	result += fmt.Sprintf("func %sGetMulti( db *sqlx.DB, key %s,searchType %s, value string) ([]%s,error){\n", info.StructName, fieldStruct, searchSruct, info.StructName)
	result += "	var result []" + info.StructName + "\n"
	result += `	statement := fmt.Sprintf("SELECT * from %s.%s where %s %s ?","` + info.DatabaseName + `",` + tableName + `,key,searchType)
	return result, db.Unsafe().Select(&result,statement,value)
}`
	return
}

func saveFunc() (result string) {
	// build func definition
	result += fmt.Sprintf("func %sSave(db *sqlx.DB, in %s) error {\n", info.StructName, info.StructName)

	var execPairs []string
	var keypairs []string
	constKey := info.StructName + "Field"
	tableName := lowerCaseFirst(info.StructName) + "TableName"
	for _, v := range info.Fields {
		cleanTag := v.Tag[1 : len(v.Tag)-1]
		tag := reflect.StructTag(cleanTag)

		if _, ok := tag.Lookup("key"); ok {
			keypairs = append(keypairs, constKey+v.Name+", "+"in."+v.Name)
			continue
		}
		if _, ok := tag.Lookup("db"); ok {
			execPairs = append(execPairs, constKey+v.Name+", "+"in."+v.Name)
		}
	}
	// build statement
	result += "\tstatement := fmt.Sprintf(\"UPDATE %s.%s SET "
	for range execPairs {
		result += "?=? "
	}
	result += "WHERE"
	for k := range keypairs {
		if k != 0 {
			result += " AND"
		}
		result += " ?=?"
	}
	result += "\", \"" + info.DatabaseName + "\", " + tableName + ")\n"
	result += "\t_,err := db.Exec(statement,\n"
	for _, v := range execPairs {
		result += "\t\t" + v + ",\n"
	}
	for _, v := range keypairs {
		result += "\t\t" + v + ",\n"
	}
	result += "\t)\n"
	result += "\treturn err\n}"
	return
}

func saveMultiFunc() (result string) {
	result += fmt.Sprintf("func %sSaveMulti(db *sqlx.DB, in []%s) error {\n", info.StructName, info.StructName)
	result += "	for _ , v := range in {\n"
	result += fmt.Sprintf("		err := %sSave(db,v)\n", info.StructName)
	result += "		if err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\treturn nil\n}"
	return
}

func newFunc() (result string) {
	// build func definition
	result += fmt.Sprintf("func %sNew(db *sqlx.DB, in %s) error {\n", info.StructName, info.StructName)
	result += "\tstatement := fmt.Sprintf(\"INSERT INTO %s.%s ("
	tableName := lowerCaseFirst(info.StructName) + "TableName"
	type pair struct {
		key   string
		value string
	}
	var pairs []pair
	for _, v := range info.Fields {
		cleanRawTag := v.Tag[1 : len(v.Tag)-1]
		tag := reflect.StructTag(cleanRawTag)
		if _, ok := tag.Lookup("db"); !ok {
			fmt.Println("No DB Found", tag)
			continue
		}
		if v, ok := tag.Lookup("key"); ok && v == "auto" {
			continue
		}
		dbtag := tag.Get("db")
		if dbtag == "" {
			continue
		}
		pairs = append(pairs, pair{key: dbtag, value: v.Name})
	}
	for k := range pairs {
		if k != 0 {
			result += ","
		}
		result += "%s"
	}
	result += ") VALUES ("
	for k := range pairs {
		if k != 0 {
			result += ","
		}
		result += "?"
	}
	result += ")\",\n\t\t\"" + info.DatabaseName + "\"," + tableName + ",\n\t\t"
	for k, v := range pairs {
		if k != 0 {
			result += ","
		}
		result += info.StructName + "Field" + v.value
	}
	result += ")\n\t_, err := db.Exec(statement,\n\t\t"
	for k, v := range pairs {
		if k != 0 {
			result += ","
		}
		result += "in." + v.value
	}
	result += ")\n\treturn err"
	result += "\n}"
	return
}

const templateText = `package {{.PackageName}}
//This Code is generated DO NOT EDIT

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type {{lowerCase .StructName}}Field string
type {{lowerCase .StructName}}SearchType string

const(
{{constants}}
)

{{getFunc}}

{{getMultiFunc}}

{{saveFunc}}

{{saveMultiFunc}}

{{newFunc}}
`
