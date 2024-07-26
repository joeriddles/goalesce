package parse

import (
	"errors"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"log"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/entity"
)

var (
	gormModelIdTag        string = "`gorm:\"primarykey\"`"
	gormModelDeletedAtTag string = "`gorm:\"index\"`"
)

type Parser interface {
	Parse(filepath string) ([]*entity.GormModelMetadata, error)
}

type parser struct {
	logger *log.Logger
	cfg    *config.Config
}

func NewParser(logger *log.Logger, cfg *config.Config) Parser {
	return &parser{
		logger: logger,
		cfg:    cfg,
	}
}

// Parse GORM model metadata from the Go file
func (p *parser) Parse(filepath string) ([]*entity.GormModelMetadata, error) {
	fset := token.NewFileSet()

	var node *ast.File
	var err error

	node, err = goparser.ParseFile(fset, filepath, nil, goparser.SkipObjectResolution) // ParseComments
	if err != nil {
		return nil, err
	}

	metadatas := []*entity.GormModelMetadata{}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			metadata, _err := p.parseGormModel(x)
			if _err != nil {
				err = _err
			}
			if metadata != nil {
				metadatas = append(metadatas, metadata)
			}
		}
		return true
	})

	return metadatas, err
}

// Parse metadata about the GORM model node
func (p *parser) parseGormModel(node *ast.TypeSpec) (*entity.GormModelMetadata, error) {
	if !p.checkIsGormModel(node) {
		msg := fmt.Sprintf("%v does not inherit from gorm.Model", node.Name.Name)
		if p.cfg.AllowCustomModels {
			p.logger.Print(msg)
		} else {
			return nil, errors.New(msg)
		}
	}

	name := node.Name.Name
	fields := p.parseGormModelFields(node)

	isGormModelEmbedded := false
	for _, field := range fields {
		isGormModelEmbedded = isGormModelEmbedded || field.IsGormModelEmbedded
		if isGormModelEmbedded {
			break
		}
	}

	metadata := &entity.GormModelMetadata{
		Name:                name,
		IsGormModelEmbedded: isGormModelEmbedded,
		Fields:              fields,
	}
	return metadata, nil
}

func (p *parser) parseGormModelFields(node *ast.TypeSpec) []*entity.GormModelField {
	fields := []*entity.GormModelField{}
	ast.Inspect(node, func(n ast.Node) bool {
		switch f := n.(type) {
		case *ast.Field:
			// embedded structs
			if len(f.Names) == 0 {
				fType := p.parseType(f.Type)
				if fType == "gorm.Model" {
					// TODO(joeriddles): handle any embedded struct
					// gormPath := reflect.TypeOf(gorm.Model{}).PkgPath()

					fields = append(
						fields,
						&entity.GormModelField{Name: "ID", Type: "uint", Tag: &gormModelIdTag, IsGormModelEmbedded: true},
						&entity.GormModelField{Name: "CreatedAt", Type: "time.Time", IsGormModelEmbedded: true},
						&entity.GormModelField{Name: "UpdatedAt", Type: "time.Time", IsGormModelEmbedded: true},
						&entity.GormModelField{Name: "DeletedAt", Type: "gorm.DeletedAt", Tag: &gormModelDeletedAtTag, IsGormModelEmbedded: true},
					)
				}
				break
			}

			fName := f.Names[0].Name // TODO(joeriddles): support multiple names
			fType := p.parseType(f.Type)

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

func (p *parser) parseType(f ast.Expr) string {
	var fType string

	switch t := f.(type) {
	case *ast.Ident:
		fType = t.Name
	case *ast.StarExpr:
		elementType := p.parseType(t.X)
		fType = fmt.Sprintf("*%v", elementType)
	case *ast.ArrayType:
		elementType := p.parseType(t.Elt)
		fType = fmt.Sprintf("[]%v", elementType)
	case *ast.SelectorExpr:
		fType = fmt.Sprintf("%v.%v", p.parseType(t.X), p.parseType(t.Sel))
	}

	return fType
}

// Check if the ast node is a GORM model
func (p *parser) checkIsGormModel(node *ast.TypeSpec) bool {
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
