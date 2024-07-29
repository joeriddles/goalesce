package entity

import (
	"fmt"
	"go/types"
	"regexp"
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

	re := regexp.MustCompile(fmt.Sprintf("%v(/\\w+?)\\.", moduleName))
	strType = re.ReplaceAllString(strType, "")

	f.Type = strType
}

func (f *GormModelField) GetType() types.Type {
	return f.t
}
