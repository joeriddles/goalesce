package parse

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/entity"
)

// Parse GORM model metadata from the Go file
func Parse(filepath string) ([]*entity.GormModelMetadata, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, nil, parser.SkipObjectResolution) // ParseComments
	if err != nil {
		return nil, err
	}

	metadatas := []*entity.GormModelMetadata{}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			metadata, err := parseGormModel(x)
			if err != nil {
				panic(err)
			}
			metadatas = append(metadatas, metadata)
		}
		return true
	})

	return metadatas, nil
}

// Parse metadata about the GORM model node
func parseGormModel(node *ast.TypeSpec) (*entity.GormModelMetadata, error) {
	if !checkIsGormModel(node) {
		return nil, errors.New("not a GORM model")
	}

	name := node.Name.Name
	fields := parseGormModelFields(node)

	metadata := &entity.GormModelMetadata{
		Name:   name,
		Fields: fields,
	}
	return metadata, nil
}

func parseGormModelFields(node *ast.TypeSpec) []*entity.GormModelField {
	fields := []*entity.GormModelField{}
	ast.Inspect(node, func(n ast.Node) bool {
		switch f := n.(type) {
		case *ast.Field:
			if len(f.Names) == 0 {
				break
			}

			fName := f.Names[0].Name // TODO(joeriddles): support multiple names
			fType := parseType(f.Type)

			var fTag *string
			if f.Tag != nil {
				fTag = &f.Tag.Value
			}

			field := &entity.GormModelField{
				Name: fName,
				Type: fType,
				Tag:  fTag,
			}
			fields = append(fields, field)
		}
		return true
	})
	return fields
}

func parseType(f ast.Expr) string {
	var fType string

	switch t := f.(type) {
	case *ast.Ident:
		fType = t.Name
	case *ast.StarExpr:
		elementType := parseType(t.X)
		fType = fmt.Sprintf("*%v", elementType)
	case *ast.ArrayType:
		elementType := parseType(t.Elt)
		fType = fmt.Sprintf("[]%v", elementType)
	}

	return fType
}

// Check if the ast node is a GORM model
func checkIsGormModel(node *ast.TypeSpec) bool {
	var isGorm = false

	ast.Inspect(node, func(n ast.Node) bool {
		switch f := n.(type) {
		case *ast.Field:
			if expr, ok := f.Type.(*ast.SelectorExpr); ok {
				xId, xIdOk := expr.X.(*ast.Ident)
				if xIdOk && xId.Name == "gorm" && expr.Sel.Name == "Model" {
					isGorm = true
				}
			}
		}
		return true
	})

	return isGorm
}
