package test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"

	"github.com/joeriddles/goalesce/pkg/v2/entity"
)

var src = "type User struct {\n" +
	"	gorm.Model\n" +
	"	Name string `gorm:\"column:name;\"`\n" +
	"}"

func ParseTestModel() *entity.GormModelMetadata {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	_, err = conf.Check("mypackage", fset, []*ast.File{f}, info)
	if err != nil {
		log.Fatal(err)
	}

}
