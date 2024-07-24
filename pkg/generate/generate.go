package generate

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/entity"
	"github.com/joeriddles/gorm-oapi-codegen/pkg/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
	"golang.org/x/tools/imports"
)

//go:embed templates
var templates embed.FS

type Generator interface {
	Generate(metadatas []*entity.GormModelMetadata) error
}

type generator struct {
	logger         *log.Logger
	outputPath     string
	moduleName     string
	modelsPkgName  string
	clearOutputDir bool

	relativePkgPath string
}

func NewGenerator(logger *log.Logger, outputPath, moduleName, modelsPkgName string, clearOutputDir bool) (Generator, error) {
	modulePath, err := utils.FindGoMod(outputPath)
	if err != nil {
		return nil, err
	}
	moduleRootPath := filepath.Dir(modulePath)
	relPath, err := filepath.Rel(moduleRootPath, outputPath)
	if err != nil {
		return nil, err
	}

	return &generator{
		logger:          logger,
		outputPath:      outputPath,
		moduleName:      moduleName,
		modelsPkgName:   modelsPkgName,
		clearOutputDir:  clearOutputDir,
		relativePkgPath: relPath,
	}, nil
}

// Generate from GORM model metadata
func (g *generator) Generate(metadatas []*entity.GormModelMetadata) error {
	t := template.New("gorm_oapi_codegen")
	if err := g.loadTemplates(templates, t); err != nil {
		return err
	}

	if g.clearOutputDir {
		if err := os.RemoveAll(g.outputPath); err != nil {
			return err
		}
	}

	if err := createDirs(
		g.outputPath,
		filepath.Join(g.outputPath, "api"),
		filepath.Join(g.outputPath, "repository"),
	); err != nil {
		return err
	}

	for _, metadata := range metadatas {
		_, err := g.generateOpenApiRoutes(t, metadata)
		if err != nil {
			return err
		}
	}

	if err := g.generateOpenApiBase(t, metadatas); err != nil {
		return err
	}

	if err := g.combineOpenApiFiles(); err != nil {
		return err
	}

	swagger, err := util.LoadSwagger(filepath.Join(g.outputPath, "openapi.yaml"))
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
	if err = os.WriteFile(filepath.Join(g.outputPath, "api", "types.gen.go"), []byte(code), 0o644); err != nil {
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
	if err = os.WriteFile(filepath.Join(g.outputPath, "api", "server_interface.gen.go"), []byte(code), 0o644); err != nil {
		return err
	}

	for _, metadata := range metadatas {
		if err := g.generateController(t, metadata); err != nil {
			return err
		}
		if err := g.generateRepository(t, metadata); err != nil {
			return err
		}
		if err := g.generateMapper(t, metadata); err != nil {
			return err
		}
	}

	if err := g.generateServer(t, metadatas); err != nil {
		return err
	}
	if err := g.generateMain(t); err != nil {
		return err
	}
	if err := g.generateMapperUtil(t); err != nil {
		return err
	}

	return nil
}

func (g *generator) combineOpenApiFiles() error {
	outputFp := filepath.Join(g.outputPath, "openapi.yaml")
	f, err := os.Create(outputFp)
	if err != nil {
		return err
	}
	defer f.Close()

	baseFp := filepath.Join(g.outputPath, "openapi_base.gen.yaml")

	args := []string{
		"bundle",
		"--output",
		outputFp,
		baseFp,
	}

	if output, err := runNpxCommand(
		"@redocly/openapi-cli@latest",
		args...,
	); err != nil {
		os.Stderr.WriteString(output)
		return err
	}

	// TODO(joeriddles): add --prune option to CLI
	entries, err := os.ReadDir(g.outputPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		filename := entry.Name()
		if strings.HasSuffix(filename, ".yaml") && filename != "openapi.yaml" {
			if err := os.Remove(filepath.Join(g.outputPath, filename)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) generateOpenApiBase(t *template.Template, metadatas []*entity.GormModelMetadata) error {
	fp := filepath.Join(g.outputPath, "openapi_base.gen.yaml")
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.ExecuteTemplate(w, "openapi_base.yaml", metadatas)
	w.Flush()

	return nil
}

func (g *generator) generateOpenApiRoutes(t *template.Template, metadata *entity.GormModelMetadata) (string, error) {
	fp := filepath.Join(g.outputPath, fmt.Sprintf("%v.gen.yaml", utils.ToSnakeCase(metadata.Name)))
	f, err := os.Create(fp)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.ExecuteTemplate(w, "openapi_controller.yaml", metadata)
	w.Flush()

	return fp, nil
}

func (g *generator) generateController(t *template.Template, metadata *entity.GormModelMetadata) error {
	fp := filepath.Join(g.outputPath, "api", fmt.Sprintf("%v_controller.gen.go", utils.ToSnakeCase(metadata.Name)))
	repositoryImportPath := filepath.Join(g.moduleName, g.relativePkgPath, "repository")
	return g.generateGo(
		t,
		fp,
		"controller.tmpl",
		map[string]interface{}{
			"repositoryImportPath": repositoryImportPath,
			"model":                metadata,
		},
	)
}

func (g *generator) generateServer(t *template.Template, metadatas []*entity.GormModelMetadata) error {
	fp := filepath.Join(g.outputPath, "api", "server.gen.go")
	return g.generateGo(
		t,
		fp,
		"server.tmpl",
		metadatas,
	)
}

func (g *generator) generateRepository(t *template.Template, metadata *entity.GormModelMetadata) error {
	fp := filepath.Join(g.outputPath, "repository", fmt.Sprintf("%v_repository.gen.go", utils.ToSnakeCase(metadata.Name)))
	return g.generateGo(
		t,
		fp,
		"repository.tmpl",
		map[string]interface{}{
			"pkg":   g.modelsPkgName,
			"model": metadata,
		},
	)
}

func (g *generator) generateMapper(t *template.Template, metadata *entity.GormModelMetadata) error {
	fp := filepath.Join(g.outputPath, "api", fmt.Sprintf("%v_mapper.gen.go", utils.ToSnakeCase(metadata.Name)))
	return g.generateGo(
		t,
		fp,
		"mapper.tmpl",
		map[string]interface{}{
			"pkg":   g.modelsPkgName,
			"model": metadata,
		},
	)
}

func (g *generator) generateMain(t *template.Template) error {
	fp := filepath.Join(g.outputPath, "main.go")
	apiImportPath := filepath.Join(g.moduleName, g.relativePkgPath, "api")
	return g.generateGo(
		t,
		fp,
		"main.tmpl",
		map[string]interface{}{
			"apiImportPath": apiImportPath,
		},
	)
}

func (g *generator) generateMapperUtil(t *template.Template) error {
	fp := filepath.Join(g.outputPath, "api", "mapper_util.gen.go")
	return g.generateGo(
		t,
		fp,
		"mapper_util.tmpl",
		map[string]interface{}{},
	)
}

// Generate formatted Go code at the filepath with the template
func (g *generator) generateGo(t *template.Template, fp string, template string, data any) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the template to in-memory buffer
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := t.ExecuteTemplate(w, template, data); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	// Format the code before saving to file
	code := b.Bytes()
	re := regexp.MustCompile("\n\\s+(\n\\s+)")
	code = re.ReplaceAll(code, []byte("$1"))

	// Format and fix missing imports
	code, err = imports.Process(fp, code, &imports.Options{})
	if err != nil {
		return err
	}

	// Write to file
	w = bufio.NewWriter(f)
	_, err = w.Write(code)
	if err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

// loadTemplates loads all of our template files into a text/template. The
// path of template is relative to the templates directory.
func (g *generator) loadTemplates(src embed.FS, t *template.Template) error {
	funcMap := template.FuncMap{
		"ToLower":        strings.ToLower,
		"ToCamelCase":    utils.ToCamelCase,
		"ToSnakeCase":    utils.ToSnakeCase,
		"ToPascalCase":   utils.ToPascalCase,
		"ToOpenApiType":  toOpenApiType,
		"MapToModelType": mapToModelType,
		"MapToApiType":   mapToApiType,
		"IsSimpleType":   isSimpleType,
		"IsComplexType":  isComplexType,
		"IsNullable":     isNullable,
		"Not":            not,
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
	Type     string
	Ref      *string
	Items    *map[string]string
	Format   *string
	Nullable bool
}

// Map the model field to a type for mapping to a model
func mapToModelType(field entity.GormModelField) string {
	result := ""
	apiField := utils.ToPascalCase(field.Name)

	switch field.Type {
	case "uint":
		result = fmt.Sprintf("%v: uint(obj.%v),", field.Name, apiField)
	case "gorm.DeletedAt":
		result = fmt.Sprintf("%v: gorm.DeletedAt{Time: *obj.%v},", field.Name, apiField)
	default:
		if isSimpleType(field.Type) {
			result = fmt.Sprintf("%v: obj.%v,", field.Name, apiField)
		} else {
			// TODO(joeriddles): map complex types
			result = ""
		}
	}

	return result
}

// Map the model field to a type for mapping to an API model
func mapToApiType(field entity.GormModelField) string {
	result := ""
	apiField := utils.ToPascalCase(field.Name)

	switch field.Type {
	case "uint":
		result = fmt.Sprintf("%v: int(model.%v),", apiField, field.Name)
	case "gorm.DeletedAt":
		result = fmt.Sprintf("%v: func() *time.Time { if model.%v.Valid { return &model.%v.Time } else { return nil } }(),", apiField, field.Name, field.Name)
	default:
		result = fmt.Sprintf("%v: model.%v,", apiField, field.Name)
	}

	return result
}

// TODO(joeriddles): Refactor this monstrosity
func toOpenApiType(t string) *OpenApiType {
	var result *OpenApiType
	nullable := false

	if isPointer := strings.HasPrefix(t, "*"); isPointer {
		t = t[1:]
		nullable = true
	} else if isArray := strings.HasPrefix(t, "[]"); isArray {
		elemType := toOpenApiType(t[2:])
		items := map[string]string{}
		if elemType.Ref != nil {
			items["$ref"] = *elemType.Ref
		} else {
			items["type"] = elemType.Type
		}

		result = &OpenApiType{
			Type:     "array",
			Items:    &items,
			Nullable: true,
		}
	}

	if result == nil {
		switch t {
		case "string":
			result = &OpenApiType{Type: "string", Nullable: nullable}
		case "time.Time":
			format := "date-time"
			result = &OpenApiType{Type: "string", Format: &format, Nullable: nullable}
		case "gorm.DeletedAt":
			format := "date-time"
			result = &OpenApiType{Type: "string", Format: &format, Nullable: true}
		case "int", "uint":
			result = &OpenApiType{Type: "integer", Nullable: nullable}
		case "int64":
			format := "int64"
			result = &OpenApiType{Type: "integer", Format: &format, Nullable: nullable}
		case "float", "float64":
			format := "float"
			result = &OpenApiType{Type: "number", Format: &format, Nullable: nullable}
		case "bool":
			result = &OpenApiType{Type: "boolean", Nullable: nullable}
		default:
			var typeRef *string = nil
			if !isSimpleType(t) {
				typeRefVal := fmt.Sprintf("./%v.gen.yaml#/components/schemas/%v", utils.ToSnakeCase(t), t)
				typeRef = &typeRefVal
				result = &OpenApiType{Type: t, Ref: typeRef, Nullable: nullable}
			} else {
				// TODO(joeriddles): panic?
				result = &OpenApiType{Type: t, Nullable: nullable}
			}
		}
	}

	return result
}

func isComplexType(t string) bool {
	if t == "" {
		return false
	}
	return strings.HasPrefix(t, "*") || !strings.HasPrefix(t, "[]") || t[0:1] != strings.ToUpper(t[0:1])
}

func isSimpleType(t string) bool {
	return !isComplexType(t)
}

func not(v bool) bool {
	return !v
}

func isNullable(t string) bool {
	openApiType := toOpenApiType(t)
	return openApiType.Nullable
}

func createDirs(paths ...string) error {
	for _, path := range paths {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			if err.Error() != "mkdir generated: file exists" {
				return err
			}
		}
	}
	return nil
}
