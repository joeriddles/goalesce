package generate

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/entity"
	"github.com/joeriddles/gorm-oapi-codegen/pkg/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

//go:embed templates
var templates embed.FS

// Generate from GORM model metadata
func Generate(metadatas []*entity.GormModelMetadata) error {
	t := template.New("gorm_oapi_codegen")
	if err := loadTemplates(templates, t); err != nil {
		return err
	}

	if err := os.RemoveAll("./generated"); err != nil {
		return err
	}

	if err := createDirs(
		"./generated",
		"./generated/api",
		"./generated/repository",
	); err != nil {
		return err
	}

	for _, metadata := range metadatas {
		_, err := generateOpenApiRoutes(t, metadata)
		if err != nil {
			return err
		}
	}

	if err := generateOpenApiBase(t, metadatas); err != nil {
		return err
	}

	if err := combineOpenApiFiles(); err != nil {
		return err
	}

	swagger, err := util.LoadSwagger("./generated/openapi.yaml")
	if err != nil {
		return err
	}

	code, err := codegen.Generate(swagger, codegen.Configuration{
		PackageName: "api",
		Generate:    codegen.GenerateOptions{Models: true},
	})
	if err != nil {
		return err
	}
	if err = os.WriteFile("./generated/api/types.gen.go", []byte(code), 0o644); err != nil {
		return err
	}

	code, err = codegen.Generate(swagger, codegen.Configuration{
		PackageName: "api",
		Generate: codegen.GenerateOptions{
			StdHTTPServer: true,
			Strict:        true,
			EmbeddedSpec:  true,
		},
	})

	if err != nil {
		return err
	}
	if err = os.WriteFile("./generated/api/server_interface.gen.go", []byte(code), 0o644); err != nil {
		return err
	}

	for _, metadata := range metadatas {
		if err := generateController(t, metadata); err != nil {
			return err
		}
		if err := generateRepository(t, metadata); err != nil {
			return err
		}
		if err := generateMapper(t, metadata); err != nil {
			return err
		}
	}

	if err := generateServer(t, metadatas); err != nil {
		return err
	}
	if err := generateMain(t); err != nil {
		return err
	}

	return nil
}

func combineOpenApiFiles() error {
	f, err := os.Create("./generated/openapi.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	args := []string{
		"bundle",
		"--output",
		"./generated/openapi.yaml",
		"./generated/openapi_base.gen.yaml",
	}

	if output, err := runNpxCommand(
		"@redocly/openapi-cli@latest",
		args...,
	); err != nil {
		os.Stderr.WriteString(output)
		return err
	}

	// TODO(joeriddles): add --prune option to CLI
	entries, err := os.ReadDir("./generated")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		filename := entry.Name()
		if strings.HasSuffix(filename, ".yaml") && filename != "openapi.yaml" {
			if err := os.Remove(fmt.Sprintf("./generated/%v", filename)); err != nil {
				return err
			}
		}
	}

	return nil

}

func generateOpenApiBase(t *template.Template, metadatas []*entity.GormModelMetadata) error {
	f, err := os.Create("generated/openapi_base.gen.yaml")
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.ExecuteTemplate(w, "openapi_base.yaml", metadatas)
	w.Flush()

	return nil
}

func generateOpenApiRoutes(t *template.Template, metadata *entity.GormModelMetadata) (string, error) {
	filename := fmt.Sprintf("./generated/%v.gen.yaml", utils.ToSnakeCase(metadata.Name))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.ExecuteTemplate(w, "openapi_controller.yaml", &metadata)
	w.Flush()

	return filename, nil
}

func generateController(t *template.Template, metadata *entity.GormModelMetadata) error {
	filepath := fmt.Sprintf("generated/api/%v_controller.gen.go", utils.ToSnakeCase(metadata.Name))
	return generateGo(
		t,
		filepath,
		"controller.tmpl",
		&metadata,
	)
}

func generateServer(t *template.Template, metadatas []*entity.GormModelMetadata) error {
	return generateGo(
		t,
		"generated/api/server.gen.go",
		"server.tmpl",
		metadatas,
	)
}

func generateRepository(t *template.Template, metadata *entity.GormModelMetadata) error {
	filepath := fmt.Sprintf("generated/repository/%v_repository.gen.go", utils.ToSnakeCase(metadata.Name))
	return generateGo(
		t,
		filepath,
		"repository.tmpl",
		&metadata,
	)
}

func generateMapper(t *template.Template, metadata *entity.GormModelMetadata) error {
	filepath := fmt.Sprintf("generated/api/%v_mapper.gen.go", utils.ToSnakeCase(metadata.Name))
	return generateGo(
		t,
		filepath,
		"mapper.tmpl",
		&metadata,
	)
}

func generateMain(t *template.Template) error {
	return generateGo(
		t,
		"generated/main.go",
		"main.tmpl",
		nil,
	)
}

// Generate formatted Go code at the filepath with the template
func generateGo(t *template.Template, filepath string, template string, data any) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.ExecuteTemplate(w, template, data)
	w.Flush()

	// outBytes, err := imports.Process(filepath, nil, &imports.Options{
	// 	FormatOnly: true,
	// })
	// if err != nil {
	// 	return err
	// }

	// w = bufio.NewWriter(f)
	// _, err = w.Write(outBytes)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// loadTemplates loads all of our template files into a text/template. The
// path of template is relative to the templates directory.
func loadTemplates(src embed.FS, t *template.Template) error {
	funcMap := template.FuncMap{
		"ToLower":       strings.ToLower,
		"ToCamelCase":   utils.ToCamelCase,
		"ToSnakeCase":   utils.ToSnakeCase,
		"ToPascalCase":  utils.ToPascalCase,
		"ToOpenApiType": toOpenApiType,
		"IsSimpleType":  isSimpleType,
		"IsId":          isId,
		"Not":           not,
	}

	return fs.WalkDir(src, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}

		buf, err := src.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file '%s': %w", path, err)
		}

		templateName := strings.TrimPrefix(path, "templates/")
		tmpl := t.New(templateName).Funcs(funcMap)
		_, err = tmpl.Parse(string(buf))
		if err != nil {
			return fmt.Errorf("parsing template '%s': %w", path, err)
		}
		return nil
	})
}

func runNpxCommand(command string, args ...string) (string, error) {
	cmd := exec.Command("npx", append([]string{command}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

type OpenApiType struct {
	Type  string
	Ref   *string
	Items *map[string]string
}

// TODO(joeriddles): Refactor this monstrosity
func toOpenApiType(t string) *OpenApiType {
	var result *OpenApiType

	if isPointer := strings.HasPrefix(t, "*"); isPointer {
		t = t[1:]
	} else if isArray := strings.HasPrefix(t, "[]"); isArray {
		elemType := toOpenApiType(t[2:])
		items := map[string]string{}
		if elemType.Ref != nil {
			items["$ref"] = *elemType.Ref
		} else {
			items["type"] = elemType.Type
		}

		result = &OpenApiType{
			Type:  "array",
			Items: &items,
		}
	}

	if result == nil {
		switch t {
		case "string":
			result = &OpenApiType{Type: "string"}
		case "int", "uint":
			result = &OpenApiType{Type: "integer"}
		default:
			var typeRef *string = nil
			if !isSimpleType(t) {
				typeRefVal := fmt.Sprintf("./%v.gen.yaml#/components/schemas/%v", utils.ToSnakeCase(t), t)
				typeRef = &typeRefVal
				result = &OpenApiType{Type: t, Ref: typeRef}
			} else {
				result = &OpenApiType{Type: t}
			}
		}
	}

	return result
}

func isSimpleType(t string) bool {
	if t == "" {
		return false // TODO(joeriddles) should this ever be empty?
	}
	return t[0:1] != strings.ToUpper(t[0:1])
}

func isId(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), "id")
}

func not(v bool) bool {
	return !v
}

func createDirs(paths ...string) error {
	for _, path := range paths {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			if err.Error() != "mkdir generated: file exists" {
				return err
			}
		}
	}
	return nil
}
