package v2

import (
	"fmt"
	"go/types"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/entity"
	"github.com/joeriddles/goalesce/pkg/parse"
	"golang.org/x/tools/go/packages"
)

const LoadAll = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax

type parser struct {
	cfg *config.Config
	pkg *packages.Package
}

func NewParser(cfg *config.Config) parse.Parser {
	return &parser{
		cfg: cfg,
	}
}

// Given a the filepath to a package with GORM models, parse all the models
func (p *parser) Parse(pkgStr string) ([]*entity.GormModelMetadata, error) {
	conf := &packages.Config{
		Mode: LoadAll,
		Dir:  p.cfg.InputFolderPath,
	}
	pkgs, err := packages.Load(conf, pkgStr)
	if err != nil {
		return nil, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("pkg %v had errors", pkgStr)
	}

	metadatas := []*entity.GormModelMetadata{}

	p.pkg = pkgs[0]
	for _, name := range p.pkg.Types.Scope().Names() {
		obj := p.pkg.Types.Scope().Lookup(name)
		if obj == nil {
			return nil, fmt.Errorf("%s.%s not found", p.pkg.Types.Path(), name)
		}

		if !obj.Exported() {
			continue
		}

		metadata, err := p.parseObject(obj.Type())
		if err != nil {
			return nil, err
		}
		if metadata != nil {
			metadatas = append(metadatas, metadata)
		}
	}

	return metadatas, nil
}

func (p *parser) parseObject(t types.Type) (*entity.GormModelMetadata, error) {
	var metadata *entity.GormModelMetadata
	var err error
	switch t := t.(type) {
	case *types.Basic:
		break
	case *types.Pointer:
		break
	case *types.Array:
		break
	case *types.Slice:
		break
	case *types.Map:
		break
	case *types.Chan:
		break
	case *types.Struct:
		metadata = p.parseStruct(t)
	case *types.Tuple:
		break
	case *types.Signature:
		break
	case *types.Named:
		metadata = p.parseNamed(t)
	case *types.Interface:
		break
	}
	return metadata, err
}

func (p *parser) parseNamed(t *types.Named) *entity.GormModelMetadata {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		metadata := p.parseStruct(u)
		metadata.Name = t.Obj().Name() // ?
		return metadata
	case *types.Map:
		return nil
	case *types.Array, *types.Slice:
		return nil
	default:
		panic("impossible")
	}
}

func (p *parser) parseStruct(t *types.Struct) *entity.GormModelMetadata {
	metadata := &entity.GormModelMetadata{
		Fields:   []*entity.GormModelField{},
		Embedded: []*entity.GormModelMetadata{},
	}

	for i := 0; i < t.NumFields(); i++ {
		field := t.Field(i)
		if !field.Exported() {
			continue
		}

		if field.Embedded() {
			fieldMetadata, _ := p.parseObject(field.Type().Underlying())
			metadata.Embedded = append(metadata.Embedded, fieldMetadata)
			continue
		}

		modelField := p.parseField(field)
		modelField.Tag = t.Tag(i)

		metadata.Fields = append(metadata.Fields, modelField)
	}

	return metadata
}

func (p *parser) parseField(field *types.Var) *entity.GormModelField {
	if field.Embedded() {
		panic("do not call parseField with embedded fields")
	}

	return &entity.GormModelField{
		Name: field.Name(),
		Type: field.Type().String(),
	}
}
