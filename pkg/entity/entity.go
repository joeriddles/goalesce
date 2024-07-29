package entity

import (
	"go/types"
	"strings"
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
	strType = strings.ReplaceAll(strType, "command-line-arguments.", "")

	if strings.HasPrefix(strType, moduleName) {
		parts := strings.Split(strType, ".")
		strType = parts[len(parts)-1]
	}

	f.Type = strType
}

func (f *GormModelField) GetType() types.Type {
	return f.t
}
