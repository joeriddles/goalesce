package api

import (
	"regexp"
	"slices"
	"strings"

	"github.com/joeriddles/goalesce/pkg/v2/builder"
	"github.com/joeriddles/goalesce/pkg/v2/entity"
	"github.com/joeriddles/goalesce/pkg/v2/generator"
	"github.com/joeriddles/goalesce/pkg/v2/utils"
)

var _ generator.ModelGenerator = new(openapiYamlControllerGenerator)

type openapiYamlControllerGenerator struct {
	model           *entity.GormModelMetadata
	yamlCodeBuilder builder.YamlCodeBuilder
}

func newOpenapiYamlControllerGenerator(
	model *entity.GormModelMetadata,
) *openapiYamlControllerGenerator {
	return &openapiYamlControllerGenerator{
		model:           model,
		yamlCodeBuilder: builder.NewYamlCodeBuilder(),
	}
}

func (o *openapiYamlControllerGenerator) DefaultOutputPath() string {
	panic("unimplemented")
}

func (o *openapiYamlControllerGenerator) EffectiveOutputPath() string {
	panic("unimplemented")
}

func (o *openapiYamlControllerGenerator) IsDisabled() bool {
	panic("unimplemented")
}

func (o *openapiYamlControllerGenerator) Generate() (string, error) {
	o.yamlCodeBuilder.DocComment("Code generated by github.com/joeriddles/goalesce DO NOT EDIT.")

	o.writePaths()
	o.writeComponents()

	code := o.yamlCodeBuilder.String()
	return code, nil
}

func (o *openapiYamlControllerGenerator) writePaths() {
	cb := o.yamlCodeBuilder
	cb.Bock("paths", func() {
		cb.Bock("/", func() {
			o.writeGetPath()
			o.writePostPath()
		})

		cb.Bock("/{id}/", func() {
			o.writeGetIDPath()
			o.writePutIDPath()
			o.writeDeleteIDPath()
		})
	})
}

func (o *openapiYamlControllerGenerator) writeGetPath() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	cb.Bock("get", func() {
		o.writeTags()
		cb.Linef("summary: Get all %vs", name)
		cb.Bock("responses", func() {
			cb.Bock(`"200"`, func() {
				cb.Line("description: Success")
				cb.Bock(("content"), func() {
					cb.Bock(("application/json"), func() {
						cb.Bock("schema", func() {
							cb.Line("type: array")
							cb.Bock("items", func() {
								cb.Linef("$ref: \"#/components/schemas/%v\"", name)
							})
						})
					})
				})

			})
		})
	})
}

func (o *openapiYamlControllerGenerator) writePostPath() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	cb.Bock("post", func() {
		o.writeTags()
		cb.Linef("summary: Create a new %v", name)
		cb.Bock("requestBody", func() {
			cb.Line("required: true")
			cb.Bock(("content"), func() {
				cb.Bock(("application/json"), func() {
					cb.Bock(("schema"), func() {
						cb.Linef("$ref: \"#/components/schemas/Create%v\"", name)
					})
				})
			})
		})
		cb.Bock("responses", func() {
			cb.Bock(`"201"`, func() {
				cb.Line("description: Created")
				cb.Bock(("content"), func() {
					cb.Bock(("application/json"), func() {
						cb.Bock(("schema"), func() {
							cb.Linef("$ref: \"#/components/schemas/%v\"", name)
						})
					})
				})
			})
			o.write400()
			o.write409()
		})
	})
}

func (o *openapiYamlControllerGenerator) writeGetIDPath() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	cb.Bock("get", func() {
		o.writeTags()
		cb.Linef("summary: Get a %v by ID", name)
		cb.Bock("parameters", func() {
			o.writeIdParameter()
		})
		cb.Bock("responses", func() {
			cb.Bock(`"200"`, func() {
				cb.Line("description: OK")
				cb.Bock(("content"), func() {
					cb.Bock(("application/json"), func() {
						cb.Bock(("schema"), func() {
							cb.Linef("$ref: \"#/components/schemas/%v\"", name)
						})
					})
				})
			})
			o.write404()
		})
	})
}

func (o *openapiYamlControllerGenerator) writePutIDPath() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	cb.Bock("put", func() {
		o.writeTags()
		cb.Linef("summary: Update a %v by ID", name)
		cb.Bock("parameters", func() {
			o.writeIdParameter()
		})
		cb.Bock("requestBody", func() {
			cb.Line("required: true")
			cb.Bock(("content"), func() {
				cb.Bock(("application/json"), func() {
					cb.Bock(("schema"), func() {
						cb.Linef("$ref: \"#/components/schemas/Update%v\"", name)
					})
				})
			})
		})
		cb.Bock("responses", func() {
			cb.Bock(`"204"`, func() {
				cb.Line("description: Updated")
			})
			o.write404()
		})
	})
}

func (o *openapiYamlControllerGenerator) writeDeleteIDPath() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	cb.Bock("delete", func() {
		o.writeTags()
		cb.Linef("summary: Delete a %v by ID", name)
		cb.Bock("parameters", func() {
			o.writeIdParameter()
		})
		cb.Bock("responses", func() {
			cb.Bock(`"204"`, func() {
				cb.Line("description: Deleted")
			})
			o.write404()
		})
	})
}

func (o *openapiYamlControllerGenerator) writeComponents() {
	cb := o.yamlCodeBuilder
	cb.Bock("components", func() {
		cb.Bock("schemas", func() {
			o.writeModelSchema()
			o.writeCreateModelSchema()
			o.writeUpdateModelSchema()
			o.writeIdSchema()
			o.writeErrorResponseSchema()
		})
		o.writeParameters()
		o.writeResponses()
	})
}

func (o *openapiYamlControllerGenerator) writeModelSchema() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	allFields := o.model.AllFields()
	cb.Bockf("%v", name)(func() {
		cb.Line("type: object")
		o.writeModelFields(allFields)
		o.writeRequiredModelFields(allFields)
	})
}

func (o *openapiYamlControllerGenerator) writeCreateModelSchema() {
	name := o.model.Name
	cb := o.yamlCodeBuilder
	// Don't include embedded fields
	fields := o.model.Fields
	cb.Bockf("Create%v", name)(func() {
		cb.Line("type: object")
		o.writeModelFields(fields)
		o.writeRequiredModelFields(fields)
	})
}

func (o *openapiYamlControllerGenerator) writeUpdateModelSchema() {
	name := o.model.Name
	cb := o.yamlCodeBuilder

	includedFields := []*entity.GormModelField{}
	for _, field := range o.model.AllFields() {
		if !o.shouldExcludeField(field) {
			includedFields = append(includedFields, field)
		}
	}

	cb.Bockf("Update%v", name)(func() {
		cb.Line("type: object")
		o.writeModelFields(includedFields)
		o.writeRequiredModelFields(includedFields)
	})
}

func (o *openapiYamlControllerGenerator) writeModelFields(fields []*entity.GormModelField) {
	o.yamlCodeBuilder.Bock("properties", func() {
		for _, field := range fields {
			o.writeModelField(field)
		}
	})
}

func (o *openapiYamlControllerGenerator) writeModelField(field *entity.GormModelField) {
	cb := o.yamlCodeBuilder
	name := utils.ToSnakeCase(field.Name)
	openApiType := utils.FieldToOpenApiType(field)

	cb.Bockf("%v", name)(func() {
		if openApiType.Ref != nil {
			cb.Linef("$ref: %v", openApiType.Ref)
		} else {
			cb.Linef("type: %v", openApiType.Type)
		}
		if openApiType.Format != nil {
			cb.Linef("format: %v", *openApiType.Format)
		}
		if openApiType.Nullable {
			cb.Line("nullable: true")
		}
		if openApiType.Items != nil {
			cb.Bock("items", func() {
				for key, value := range *openApiType.Items {
					cb.Linef("%v: \"%v\"", key, value)
				}
			})
		}
	})
}

func (o *openapiYamlControllerGenerator) writeRequiredModelFields(fields []*entity.GormModelField) {
	cb := o.yamlCodeBuilder
	cb.Bock("required", func() {
		for _, field := range fields {
			openApiType := utils.FieldToOpenApiType(field)
			if !openApiType.Nullable {
				fieldName := utils.ToSnakeCase(field.Name)
				cb.List(fieldName, func() {})
			}
		}
	})
}

func (o *openapiYamlControllerGenerator) writeParameters() {
	cb := o.yamlCodeBuilder
	cb.Bock("parameters", func() {
		cb.Bock("IdPath", func() {
			cb.Line("name: id")
			cb.Line("in: path")
			cb.Line("required: true")
			cb.Bock("schema", func() {
				cb.Line("$ref: \"#/components/schemas/id\"")
			})
		})
	})
}

func (o *openapiYamlControllerGenerator) writeResponses() {
	responses := map[string]string{
		"BadRequest":   "Bad request - Contents of the request are unexpected",
		"Unauthorized": "Unauthorized - Invalid app check token, bearer token, or scope",
		"Forbidden":    "Forbidden - No permission to access the resource",
		"NotFound":     "Not Found - Specified resource could not be located",
		"Conflict":     "Conflict - Operation would result in resource conflicts",
	}
	cb := o.yamlCodeBuilder
	cb.Bock("responses", func() {
		for response, description := range responses {
			cb.Bock(response, func() {
				cb.Linef("description: \"%v\"", description)
				cb.Bock("content", func() {
					cb.Bock("application/json", func() {
						cb.Bock("schema", func() {
							cb.Line("$ref: \"#/components/schemas/ErrorResponse\"")
						})
					})
				})
			})
		}
	})

}

func (o *openapiYamlControllerGenerator) writeTags() {
	snakeCaseName := utils.ToSnakeCase(o.model.Name)
	o.yamlCodeBuilder.Bock("tags", func() {
		o.yamlCodeBuilder.Linef(`- "%v"`, snakeCaseName)
	})
}

func (o *openapiYamlControllerGenerator) writeIdSchema() {
	cb := o.yamlCodeBuilder
	cb.Bock("id", func() {
		cb.Line("type: integer")
		cb.Line("format: int64")
		cb.Line("description: A unique ID to represent a resource")
		cb.Line("minimum: 0")
	})
}

func (o *openapiYamlControllerGenerator) writeErrorResponseSchema() {
	cb := o.yamlCodeBuilder
	cb.Bock("ErrorResponse", func() {
		cb.Line("type: object")
		cb.Bock("properties", func() {
			cb.Bock("code", func() {
				cb.Line("type: string")
				cb.Line("description: The error code's unique identifier")
			})
			cb.Bock("message", func() {
				cb.Line("type: string")
				cb.Line("description: The error code's detailed message providing information about itself")
			})
		})
		cb.Bock("required", func() {
			cb.List("code", func() {})
			cb.List("message", func() {})
		})

	})

}

func (o *openapiYamlControllerGenerator) writeIdParameter() {
	cb := o.yamlCodeBuilder
	cb.List("name: id", func() {
		cb.Line("in: path")
		cb.Line("required: true")
		cb.Bock("schema", func() {
			cb.Line("$ref: \"#/components/schemas/id\"")
		})
	})
}

func (o *openapiYamlControllerGenerator) write400() {
	o.yamlCodeBuilder.Bock(`"400"`, func() {
		o.yamlCodeBuilder.Line("$ref: \"#/components/responses/BadRequest\"")
	})
}

func (o *openapiYamlControllerGenerator) write404() {
	o.yamlCodeBuilder.Bock(`"404"`, func() {
		o.yamlCodeBuilder.Line("$ref: \"#/components/responses/NotFound\"")
	})
}

func (o *openapiYamlControllerGenerator) write409() {
	o.yamlCodeBuilder.Bock(`"409"`, func() {
		o.yamlCodeBuilder.Line("$ref: \"#/components/responses/Conflict\"")
	})
}

var primaryKeyRegex *regexp.Regexp = regexp.MustCompile("gorm:\"(.*?)\"")

// Whether the field should be excluded from create and update operations
func (o *openapiYamlControllerGenerator) shouldExcludeField(field *entity.GormModelField) bool {
	if field.Parent != nil {
		if field.Parent.Pkg == "gorm" && field.Parent.Name == "Model" {
			return true
		}
	}

	match := primaryKeyRegex.FindStringSubmatch(field.Tag)
	if len(match) == 0 {
		return false
	}

	gormParts := strings.Split(match[1], ";")
	isPrimaryKey := slices.Contains(gormParts, "primaryKey")
	isAutoCreateTime := slices.ContainsFunc(gormParts, func(p string) bool { return strings.HasPrefix(p, "autoCreateTime") })
	isAutoUpdateTime := slices.ContainsFunc(gormParts, func(p string) bool { return strings.HasPrefix(p, "isAutoUpdateTime") })
	return isPrimaryKey || isAutoCreateTime || isAutoUpdateTime
}
