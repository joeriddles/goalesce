package parse

import (
	"fmt"
	"go/types"
	"log"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/entity"
	"golang.org/x/tools/go/packages"
)

type Parser interface {
	Parse(filepath string) ([]*entity.GormModelMetadata, error)
}

const LoadAll = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax

type parser struct {
	cfg    *config.Config
	logger *log.Logger
	pkg    *packages.Package
}

func NewParser(
	logger *log.Logger,
	cfg *config.Config,
) Parser {
	return &parser{
		logger: logger,
		cfg:    cfg,
	}
}

// Given a the filepath to a package with GORM models, parse all the models
func (p *parser) Parse(pkgStr string) ([]*entity.GormModelMetadata, error) {
	conf := &packages.Config{
		Mode: LoadAll,
		Dir:  pkgStr,
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
		metadata.Name = t.Obj().Name()
		return metadata
	case *types.Map:
		return nil
	case *types.Array, *types.Slice:
		return nil
	case *types.Interface:
		return nil
	case *types.Signature:
		return nil
	default:
		panic(fmt.Sprintf("parseNamed is impossible: %v", t.Obj().Name()))
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

		modelField := &entity.GormModelField{}
		modelField.Name = field.Name()
		modelField.WithType(field.Type(), p.cfg.ModuleName)
		modelField.Tag = t.Tag(i)
		metadata.Fields = append(metadata.Fields, modelField)
	}

	return metadata
}
