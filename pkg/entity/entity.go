package entity

import (
	"go/types"

	"github.com/joeriddles/goalesce/pkg/utils"
)

type GormModelMetadata struct {
	Name     string
	Fields   []*GormModelField
	Embedded []*GormModelMetadata
}

type GormModelField struct {
	Name string
	Type string
	Tag  string

	t types.Type
}

func (f *GormModelField) WithType(t types.Type, moduleName string) {
	f.t = t
	strType := t.String()

	// Ignore special "command-line-arguments" package when no explict package is specified
	strType = utils.StripModulePackage(strType, "command-line-arguments")

	strType = utils.StripModulePackage(strType, moduleName)

	f.Type = strType
}

func (f *GormModelField) GetType() types.Type {
	return f.t
}
