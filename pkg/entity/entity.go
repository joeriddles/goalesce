package entity

import (
	"go/types"
	"strings"

	"github.com/joeriddles/goalesce/pkg/utils"
)

type GormModelMetadata struct {
	Name     string
	Fields   []*GormModelField
	Embedded []*GormModelMetadata
	IsApi    bool

	t types.Type
}

func (m *GormModelMetadata) WithType(t types.Type) {
	m.t = t
}

func (m *GormModelMetadata) GetType() types.Type {
	return m.t
}

func (m *GormModelMetadata) AllFields() []*GormModelField {
	fields := []*GormModelField{}
	fields = append(fields, m.Fields...)
	for _, embedded := range m.Embedded {
		fields = append(fields, embedded.AllFields()...)
	}
	return fields
}

func (m *GormModelMetadata) GetField(name string) *GormModelField {
	field, err := utils.First(m.AllFields(), func(field *GormModelField) bool {
		return field.Name == name
	})
	if err != nil {
		// 100% a developer error
		panic(err)
	}
	return field
}

type GormModelField struct {
	Name        string
	Type        string
	Tag         string
	OpenApiType string

	MapFunc    *string
	MapApiFunc *string

	Parent *GormModelMetadata `json:"-"`

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

// Get the type for use in a Go type declaration
//
// Example: github.com/shopspring/decimal.Decimal -> decimal.Decimal
func (f *GormModelField) GetGoType() string {
	if strings.Contains(f.Type, "/") {
		split := strings.Split(f.Type, "/")
		return split[len(split)-1]
	}
	return f.Type
}
