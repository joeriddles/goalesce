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

func (m *GormModelMetadata) AllFields() []*GormModelField {
	fields := []*GormModelField{}
	for _, field := range m.Fields {
		fields = append(fields, field)
	}
	for _, embedded := range m.Embedded {
		for _, field := range embedded.AllFields() {
			fields = append(fields, field)
		}
	}
	return fields
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
