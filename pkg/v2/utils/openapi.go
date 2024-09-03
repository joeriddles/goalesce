package utils

import (
	"fmt"
	"strings"
)

// Metadata for generating an OpenAPI field
type OpenApiType struct {
	Type     string
	Ref      *string
	Items    *map[string]string
	Format   *string
	Nullable bool
}

// TODO(joeriddles): Refactor this monstrosity
func ToOpenApiType(typ string) *OpenApiType {
	if isPointer := strings.HasPrefix(typ, "*"); isPointer {
		typ = typ[1:]
		result := ToOpenApiType(typ)
		result.Nullable = true
		return result
	}

	if isArray := strings.HasPrefix(typ, "[]"); isArray {
		typ = typ[2:]
		elemType := ToOpenApiType(typ)
		elemType.Nullable = false
		items := map[string]string{}
		if elemType.Ref != nil {
			items["$ref"] = *elemType.Ref
		} else {
			items["type"] = elemType.Type
		}

		return &OpenApiType{
			Type:     "array",
			Items:    &items,
			Nullable: true,
		}
	}

	var result *OpenApiType
	switch typ {
	case "string":
		result = &OpenApiType{Type: "string"}
	case "time.Time":
		format := "date-time"
		result = &OpenApiType{Type: "string", Format: &format}
	case "gorm.io/gorm.DeletedAt":
		format := "date-time"
		result = &OpenApiType{Type: "string", Format: &format, Nullable: true}
	case "int", "uint":
		result = &OpenApiType{Type: "integer"}
	case "int64":
		format := "int64"
		result = &OpenApiType{Type: "integer", Format: &format}
	case "float", "float64":
		format := "float"
		result = &OpenApiType{Type: "number", Format: &format}
	case "bool":
		result = &OpenApiType{Type: "boolean"}
	default:
		var typeRef *string = nil
		if !IsSimpleType(typ) {
			typeRefVal := fmt.Sprintf("./%v.gen.yaml#/components/schemas/%v", ToSnakeCase(typ), typ)
			typeRef = &typeRefVal
			result = &OpenApiType{Type: typ, Ref: typeRef}
		} else {
			// TODO(joeriddles): panic?
			result = &OpenApiType{Type: typ}
		}
	}

	return result
}
