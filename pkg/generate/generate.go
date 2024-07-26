package generate

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/entity"
	"github.com/joeriddles/goalesce/pkg/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
	"golang.org/x/tools/imports"
	"gopkg.in/yaml.v2"
)

//go:embed templates
var templates embed.FS

type Generator interface {
	Generate(metadatas []*entity.GormModelMetadata) error
}

type generator struct {
	logger            *log.Logger
	cfg               *config.Config
	relativePkgPath   string
	typesPackage      *string
	repositoryPackage string
}

func NewGenerator(logger *log.Logger, cfg *config.Config) (Generator, error) {
	modulePath, err := utils.FindGoMod(cfg.OutputFile)
	if err != nil {
		return nil, err
	}
	moduleRootPath := filepath.Dir(modulePath)
	relPath, err := filepath.Rel(moduleRootPath, cfg.OutputFile)
	if err != nil {
		return nil, err
	}

	var typesPackage *string = nil
	defaultOutputFile := filepath.Join(cfg.OutputFile, "api", "types.gen.go")
	if cfg.TypesCodegen.OutputFile != defaultOutputFile {
		relPkg, err := filepath.Rel(moduleRootPath, cfg.TypesCodegen.OutputFile)
		if err != nil {
			return nil, err
		}
		pkg := filepath.Join(cfg.ModuleName, relPkg)
		pkg = filepath.Dir(pkg) // remove filename
		typesPackage = &pkg
	}

	var repositoryPackage string = filepath.Join(cfg.ModuleName, relPath, "repository")
	defaultRepositoryOutputFile := filepath.Join(cfg.OutputFile, "repository")
	if cfg.RepositoryConfiguration.OutputFile != defaultRepositoryOutputFile {
		relPkg, err := filepath.Rel(moduleRootPath, cfg.RepositoryConfiguration.OutputFile)
		if err != nil {
			return nil, err
		}
		pkg := filepath.Join(cfg.ModuleName, relPkg)
		repositoryPackage = filepath.Dir(pkg) // remove filename
	}

	return &generator{
		logger:            logger,
		cfg:               cfg,
		relativePkgPath:   relPath,
		typesPackage:      typesPackage,
		repositoryPackage: repositoryPackage,
	}, nil
}

// Generate from GORM model metadata
func (g *generator) Generate(metadatas []*entity.GormModelMetadata) error {
	t := template.New("gorm_oapi_codegen")
	if err := g.loadTemplates(templates, t); err != nil {
		return err
	}

	if g.cfg.ClearOutputDir {
		if err := os.RemoveAll(g.cfg.OutputFile); err != nil {
			return err
		}
	}

	if err := createDirs(
		g.cfg.OutputFile,
		filepath.Join(g.cfg.OutputFile, "api"),
		filepath.Join(g.cfg.RepositoryConfiguration.OutputFile),
	); err != nil {
		return err
	}

	for _, metadata := range metadatas {
		_, err := g.generateOpenApiRoutes(t, metadata)
		if err != nil {
			return err
		}
	}

	if err := g.generateOpenApiBase(metadatas); err != nil {
		return err
	}

	if err := g.combineOpenApiFiles(); err != nil {
		return err
	}

	swagger, err := util.LoadSwagger(filepath.Join(g.cfg.OutputFile, "openapi.yaml"))
	if err != nil {
		return err
	}

	code, err := codegen.Generate(swagger, g.cfg.TypesCodegen.Configuration)
	if err != nil {
		return err
	}
	if err = os.WriteFile(g.cfg.TypesCodegen.OutputFile, []byte(code), 0o644); err != nil {
		return err
	}

	code, err = codegen.Generate(swagger, g.cfg.ServerCodegen.Configuration)
	if err != nil {
		return err
	}
	if err = os.WriteFile(g.cfg.ServerCodegen.OutputFile, []byte(code), 0o644); err != nil {
		return err
	}

	for _, metadata := range metadatas {
		if slices.Contains(g.cfg.ExcludeModels, metadata.Name) {
			continue
		}

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

	filteredMetadatas := []*entity.GormModelMetadata{}
	for _, metadata := range metadatas {
		if !slices.Contains(g.cfg.ExcludeModels, metadata.Name) {
			filteredMetadatas = append(filteredMetadatas, metadata)
		}
	}
	if err := g.generateServer(t, filteredMetadatas); err != nil {
		return err
	}

	if g.cfg.GenerateMain {
		if err := g.generateMain(t); err != nil {
			return err
		}
	}

	if err := g.generateMapperUtil(t); err != nil {
		return err
	}

	return nil
}

func (g *generator) combineOpenApiFiles() error {
	outputFp := filepath.Join(g.cfg.OutputFile, "openapi.yaml")
	f, err := os.Create(outputFp)
	if err != nil {
		return err
	}
	defer f.Close()

	baseFp := filepath.Join(g.cfg.OutputFile, "openapi_base.gen.yaml")

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

	if g.cfg.PruneYaml {
		entries, err := os.ReadDir(g.cfg.OutputFile)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			filename := entry.Name()
			if strings.HasSuffix(filename, ".yaml") && filename != "openapi.yaml" {
				if err := os.Remove(filepath.Join(g.cfg.OutputFile, filename)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (g *generator) generateOpenApiBase(metadatas []*entity.GormModelMetadata) error {
	fp := filepath.Join(g.cfg.OutputFile, "openapi_base.gen.yaml")
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	var doc *openapi3.T

	if g.cfg.OpenApiFile != "" {
		loader := openapi3.NewLoader()
		doc, err = loader.LoadFromFile(g.cfg.OpenApiFile)
		if err != nil {
			return err
		}
	} else {
		doc = &openapi3.T{
			OpenAPI: "3.0.0",
			Info: &openapi3.Info{
				Version: "1.0.0",
				Title:   "Generated API",
			},
			Paths: openapi3.NewPaths(),
		}
	}

	for _, metadata := range metadatas {
		if slices.Contains(g.cfg.ExcludeModels, metadata.Name) {
			continue
		}

		doc.Paths.Set(fmt.Sprintf("/%v/", utils.ToHtmlCase(metadata.Name)), &openapi3.PathItem{
			Ref: fmt.Sprintf("./%v.gen.yaml#/paths/~1", utils.ToSnakeCase(metadata.Name)),
		})
		doc.Paths.Set(fmt.Sprintf("/%v/{id}/", utils.ToHtmlCase(metadata.Name)), &openapi3.PathItem{
			Ref: fmt.Sprintf("./%v.gen.yaml#/paths/~1%%7Bid%%7D~1", utils.ToSnakeCase(metadata.Name)),
		})
	}

	if err := doc.Paths.Validate(context.Background()); err != nil {
		return err
	}

	yamlContent, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	if _, err = f.Write(yamlContent); err != nil {
		return err
	}

	return nil
}

func (g *generator) generateOpenApiRoutes(t *template.Template, metadata *entity.GormModelMetadata) (string, error) {
	fp := filepath.Join(g.cfg.OutputFile, fmt.Sprintf("%v.gen.yaml", utils.ToSnakeCase(metadata.Name)))
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
	fp := filepath.Join(g.cfg.OutputFile, "api", fmt.Sprintf("%v_controller.gen.go", utils.ToSnakeCase(metadata.Name)))

	template := "controller.tmpl"
	if g.cfg.ServerCodegen.Generate.EchoServer {
		template = "echo_controller.tmpl"
	}

	return g.generateGo(
		t,
		fp,
		template,
		map[string]interface{}{
			"package":              g.cfg.ServerCodegen.PackageName,
			"typesPackage":         g.typesPackage,
			"repositoryImportPath": g.repositoryPackage,
			"model":                metadata,
		},
	)
}

func (g *generator) generateServer(t *template.Template, metadatas []*entity.GormModelMetadata) error {
	fp := filepath.Join(g.cfg.OutputFile, "api", "server.gen.go")

	template := "server.tmpl"
	if g.cfg.ServerCodegen.Generate.EchoServer {
		template = "echo_server.tmpl"
	}

	return g.generateGo(
		t,
		fp,
		template,
		map[string]interface{}{
			"package":      g.cfg.ServerCodegen.PackageName,
			"typesPackage": g.typesPackage,
			"metadatas":    metadatas,
		},
	)
}

func (g *generator) generateRepository(t *template.Template, metadata *entity.GormModelMetadata) error {
	filename := fmt.Sprintf("%v_repository.gen.go", utils.ToSnakeCase(metadata.Name))
	fp := filepath.Join(g.cfg.RepositoryConfiguration.OutputFile, filename)
	return g.generateGo(
		t,
		fp,
		"repository.tmpl",
		map[string]interface{}{
			"pkg":   g.cfg.ModelsPkg,
			"model": metadata,
		},
	)
}

func (g *generator) generateMapper(t *template.Template, metadata *entity.GormModelMetadata) error {
	fp := filepath.Join(g.cfg.OutputFile, "api", fmt.Sprintf("%v_mapper.gen.go", utils.ToSnakeCase(metadata.Name)))
	return g.generateGo(
		t,
		fp,
		"mapper.tmpl",
		map[string]interface{}{
			"package":      g.cfg.ServerCodegen.PackageName,
			"typesPackage": g.typesPackage,
			"pkg":          g.cfg.ModelsPkg,
			"model":        metadata,
		},
	)
}

func (g *generator) generateMain(t *template.Template) error {
	fp := filepath.Join(g.cfg.OutputFile, "main.go")
	apiImportPath := filepath.Join(g.cfg.ModuleName, g.relativePkgPath, "api")
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
	fp := filepath.Join(g.cfg.OutputFile, "api", "mapper_util.gen.go")
	return g.generateGo(
		t,
		fp,
		"mapper_util.tmpl",
		map[string]interface{}{
			"package": g.cfg.ServerCodegen.PackageName,
		},
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
	getTypesNamespace := func() string {
		if g.typesPackage == nil {
			return ""
		}
		return "types."
	}

	funcMap := template.FuncMap{
		"ToLower":        strings.ToLower,
		"ToCamelCase":    utils.ToCamelCase,
		"ToSnakeCase":    utils.ToSnakeCase,
		"ToHtmlCase":     utils.ToHtmlCase,
		"ToPascalCase":   utils.ToPascalCase,
		"ToOpenApiType":  toOpenApiType,
		"MapToModelType": mapToModelType,
		"MapToApiType":   mapToApiType,
		"IsSimpleType":   isSimpleType,
		"IsComplexType":  isComplexType,
		"IsNullable":     isNullable,
		"Not":            not,
		"Types":          getTypesNamespace,
	}

	err := fs.WalkDir(src, "templates", func(path string, d fs.DirEntry, err error) error {
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
	if err != nil {
		return err
	}

	if g.cfg.RepositoryConfiguration.Template != nil {
		path := *g.cfg.RepositoryConfiguration.Template
		buf, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file '%s': %w", path, err)
		}

		templateName := "repository.tmpl"
		tmpl := t.New(templateName).Funcs(funcMap)
		_, err = tmpl.Parse(string(buf))
		if err != nil {
			return fmt.Errorf("parsing template '%s': %w", path, err)
		}
	}

	return nil
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
				return fmt.Errorf("could not create path %v: %v", path, err)
			}
		}
	}
	return nil
}
