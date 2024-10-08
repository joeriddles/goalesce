package generate

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"go/types"
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
	"github.com/joeriddles/goalesce/pkg/convert"
	"github.com/joeriddles/goalesce/pkg/entity"
	"github.com/joeriddles/goalesce/pkg/parse"
	"github.com/joeriddles/goalesce/pkg/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	codegen_util "github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
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
	templates         *template.Template
	relativePkgPath   string
	typesPackage      *string
	repositoryPackage string
}

func NewGenerator(logger *log.Logger, cfg *config.Config) (Generator, error) {
	modulePath, err := utils.FindGoMod(cfg.OutputFile, cfg.ModuleName)
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
		pkg = filepathToPackage(pkg)
		typesPackage = &pkg
	}

	var repositoryPackage string = filepath.Join(cfg.ModuleName, relPath, "repository")
	defaultRepositoryOutputFile := filepath.Join(cfg.OutputFile, "repository")
	if cfg.RepositoryConfiguration.OutputFile != defaultRepositoryOutputFile {
		relPkg, err := filepath.Rel(moduleRootPath, cfg.RepositoryConfiguration.OutputFile)
		if err != nil {
			return nil, err
		}
		repositoryPackage = filepath.Join(cfg.ModuleName, relPkg)
	}

	// Normalize filepath to Go package
	repositoryPackage = filepathToPackage(repositoryPackage)

	g := &generator{
		logger:            logger,
		cfg:               cfg,
		relativePkgPath:   relPath,
		typesPackage:      typesPackage,
		repositoryPackage: repositoryPackage,
	}

	t := template.New("gorm_oapi_codegen")
	if err := g.loadTemplates(templates, t); err != nil {
		return nil, err
	}
	g.templates = t

	return g, nil
}

// Generate from GORM model metadata
func (g *generator) Generate(metadatas []*entity.GormModelMetadata) error {
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

	if err := g.generateOpenApiYaml(metadatas); err != nil {
		return err
	}

	if err := g.runCodegenTool(); err != nil {
		return err
	}

	// parse generated API types
	parser := parse.NewParser(g.logger, g.cfg)
	apiMetadatas, err := parser.Parse(filepath.Dir(g.cfg.TypesCodegen.OutputFile))
	if err != nil {
		return err
	}

	for _, metadata := range metadatas {
		apiMetadata, err := utils.First(apiMetadatas, func(m *entity.GormModelMetadata) bool {
			return m.Name == metadata.Name
		})
		if err != nil {
			g.logger.Printf("could not find apiMetadata for %v", metadata.Name)
			continue
		}
		apiMetadata.IsApi = true

		for _, apiField := range apiMetadata.AllFields() {
			field := metadata.GetField(apiField.Name)
			apiField.MapFunc = field.MapApiFunc
			apiField.MapApiFunc = field.MapFunc
		}

		createStr := fmt.Sprintf("Create%v", metadata.Name)
		createApiMetadata, _ := utils.First(apiMetadatas, func(m *entity.GormModelMetadata) bool {
			return m.Name == createStr
		})
		updateStr := fmt.Sprintf("Update%v", metadata.Name)
		updateApiMetadata, _ := utils.First(apiMetadatas, func(m *entity.GormModelMetadata) bool {
			return m.Name == updateStr
		})

		if err := g.generateMapper(metadata, apiMetadata); err != nil {
			return err
		}

		// Don't generate anything but the mapper for excluded models
		if slices.Contains(g.cfg.ExcludeModels, metadata.Name) {
			continue
		}

		for _, createApiField := range createApiMetadata.AllFields() {
			field := metadata.GetField(createApiField.Name)
			createApiField.MapFunc = field.MapApiFunc
			createApiField.MapApiFunc = field.MapFunc
		}
		for _, updateApiField := range updateApiMetadata.AllFields() {
			field := metadata.GetField(updateApiField.Name)
			updateApiField.MapFunc = field.MapApiFunc
			updateApiField.MapApiFunc = field.MapFunc
		}

		if err := g.generateRepository(metadata); err != nil {
			return err
		}

		getFilterStr := fmt.Sprintf("Get%vParams", metadata.Name)
		getFilterMetadata, _ := utils.First(apiMetadatas, func(m *entity.GormModelMetadata) bool {
			return m.Name == getFilterStr
		})

		for _, filterField := range getFilterMetadata.Fields {
			modelField, _ := utils.First(metadata.Fields, func(f *entity.GormModelField) bool {
				return f.Name == filterField.Name
			})
			filterField.MapFunc = modelField.MapFunc
			filterField.MapApiFunc = modelField.MapApiFunc
		}
		if err := g.generateController(metadata, createApiMetadata, updateApiMetadata, getFilterMetadata); err != nil {
			return err
		}
	}

	filteredMetadatas := []*entity.GormModelMetadata{}
	for _, metadata := range metadatas {
		if !slices.Contains(g.cfg.ExcludeModels, metadata.Name) {
			filteredMetadatas = append(filteredMetadatas, metadata)
		}
	}

	if g.cfg.GenerateServer {
		if err := g.generateServer(filteredMetadatas); err != nil {
			return err
		}
	}

	if g.cfg.GenerateMain {
		if err := g.generateMain(); err != nil {
			return err
		}
	}

	if err := g.generateMapperUtil(); err != nil {
		return err
	}

	return nil
}

func (g *generator) generateOpenApiYaml(metadatas []*entity.GormModelMetadata) error {
	for _, metadata := range metadatas {
		_, err := g.generateOpenApiRoutes(metadata)
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
	return nil
}

func (g *generator) runCodegenTool() error {
	g.cfg.TypesCodegen.Configuration.OutputOptions.NameNormalizer = "ToCamelCaseWithInitialisms"
	g.cfg.ServerCodegen.Configuration.OutputOptions.NameNormalizer = "ToCamelCaseWithInitialisms"

	swagger, err := codegen_util.LoadSwagger(filepath.Join(g.cfg.OutputFile, "openapi.yaml"))
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
		doc.Paths.Set(fmt.Sprintf("/%v/batch/", utils.ToHtmlCase(metadata.Name)), &openapi3.PathItem{
			Ref: fmt.Sprintf("./%v.gen.yaml#/paths/~1batch~1", utils.ToSnakeCase(metadata.Name)),
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

func (g *generator) generateOpenApiRoutes(metadata *entity.GormModelMetadata) (string, error) {
	fp := filepath.Join(g.cfg.OutputFile, fmt.Sprintf("%v.gen.yaml", utils.ToSnakeCase(metadata.Name)))
	f, err := os.Create(fp)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	g.templates.ExecuteTemplate(w, "openapi_controller.yaml", metadata)
	w.Flush()

	return fp, nil
}

func (g *generator) generateController(
	metadata *entity.GormModelMetadata,
	createApiMetadata *entity.GormModelMetadata,
	updateApiMetadata *entity.GormModelMetadata,
	filterMetadata *entity.GormModelMetadata,
) error {
	fp := filepath.Join(g.cfg.OutputFile, "api", fmt.Sprintf("%v_controller.gen.go", utils.ToSnakeCase(metadata.Name)))

	convertToFilter := func(field *entity.GormModelField) string {
		return convert.ConvertFieldNamed(g.templates, field, metadata, "request.Params", "filters")
	}
	g.templates.Funcs(template.FuncMap{
		"ConvertToFilter": convertToFilter,
	})

	template := "controller.tmpl"
	if g.cfg.ServerCodegen.Generate.EchoServer {
		template = "echo_controller.tmpl"
	}

	return g.generateGo(
		fp,
		template,
		map[string]interface{}{
			"package":              g.cfg.ServerCodegen.PackageName,
			"queryPackage":         g.cfg.QueryPkg,
			"typesPackage":         g.typesPackage,
			"repositoryImportPath": g.repositoryPackage,
			"model":                metadata,
			"createApi":            createApiMetadata,
			"updateApi":            updateApiMetadata,
			"filterMetadata":       filterMetadata,
		},
	)
}

func (g *generator) generateServer(metadatas []*entity.GormModelMetadata) error {
	fp := filepath.Join(g.cfg.OutputFile, "api", "server.gen.go")

	template := "server.tmpl"
	if g.cfg.ServerCodegen.Generate.EchoServer {
		template = "echo_server.tmpl"
	}

	return g.generateGo(
		fp,
		template,
		map[string]interface{}{
			"package":      g.cfg.ServerCodegen.PackageName,
			"queryPackage": g.cfg.QueryPkg,
			"typesPackage": g.typesPackage,
			"metadatas":    metadatas,
		},
	)
}

func (g *generator) generateRepository(metadata *entity.GormModelMetadata) error {
	filename := fmt.Sprintf("%v_repository.gen.go", utils.ToSnakeCase(metadata.Name))
	fp := filepath.Join(g.cfg.RepositoryConfiguration.OutputFile, filename)
	return g.generateGo(
		fp,
		"repository.tmpl",
		map[string]interface{}{
			"pkg":      g.cfg.ModelsPkg,
			"queryPkg": g.cfg.QueryPkg,
			"model":    metadata,
		},
	)
}

func (g *generator) generateMapper(
	metadata *entity.GormModelMetadata,
	apiMetadata *entity.GormModelMetadata,
) error {
	fp := filepath.Join(g.cfg.OutputFile, "api", fmt.Sprintf("%v_mapper.gen.go", utils.ToSnakeCase(metadata.Name)))
	apiFp := filepath.Join(g.cfg.OutputFile, "api", fmt.Sprintf("%v_api_mapper.gen.go", utils.ToSnakeCase(metadata.Name)))

	convertToModel := func(field *entity.GormModelField) string {
		return convert.ConvertField(g.templates, field, metadata)
	}
	convertToApi := func(field *entity.GormModelField) string {
		return convert.ConvertField(g.templates, field, apiMetadata)
	}
	g.templates.Funcs(template.FuncMap{
		"ConvertToModel": convertToModel,
		"ConvertToApi":   convertToApi,
	})

	errs := []error{}

	err := g.generateGo(
		fp,
		"model_mapper.tmpl",
		map[string]interface{}{
			"package":      g.cfg.ServerCodegen.PackageName,
			"typesPackage": g.typesPackage,
			"pkg":          g.cfg.ModelsPkg,
			"Name":         metadata.Name,
			"src":          apiMetadata,
			"dst":          metadata,
		},
	)
	if err != nil {
		errs = append(errs, err)
	}

	err = g.generateGo(
		apiFp,
		"api_mapper.tmpl",
		map[string]interface{}{
			"package":      g.cfg.ServerCodegen.PackageName,
			"typesPackage": g.typesPackage,
			"pkg":          g.cfg.ModelsPkg,
			"Name":         apiMetadata.Name + "Api",
			"src":          metadata,
			"dst":          apiMetadata,
		},
	)
	if err != nil {
		errs = append(errs, err)
	}
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (g *generator) generateMain() error {
	fp := filepath.Join(g.cfg.OutputFile, "main.go")
	apiImportPath := filepath.Join(g.cfg.ModuleName, g.relativePkgPath, "api")
	return g.generateGo(
		fp,
		"main.tmpl",
		map[string]interface{}{
			"apiImportPath": apiImportPath,
		},
	)
}

func (g *generator) generateMapperUtil() error {
	fp := filepath.Join(g.cfg.OutputFile, "api", "mapper_util.gen.go")
	return g.generateGo(
		fp,
		"mapper_util.tmpl",
		map[string]interface{}{
			"package": g.cfg.ServerCodegen.PackageName,
		},
	)
}

// Generate formatted Go code at the filepath with the template
func (g *generator) generateGo(fp string, template string, data any) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the template to in-memory buffer
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := g.templates.ExecuteTemplate(w, template, data); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	// Format code and fix missing imports
	unformattedCode := b.Bytes()
	code, err := imports.Process(fp, unformattedCode, &imports.Options{})
	if err != nil {
		fmt.Printf("Failed to format code:\n%v", string(unformattedCode))
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

	// Format imports with goimports tool
	cmd := exec.Command("goimports", "-w", fp)
	_, err = cmd.CombinedOutput()
	if err != nil {
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
		"ToLower":            strings.ToLower,
		"ToCamelCase":        utils.ToCamelCase,
		"ToSnakeCase":        utils.ToSnakeCase,
		"ToHtmlCase":         utils.ToHtmlCase,
		"ToPascalCase":       utils.ToPascalCase,
		"ShouldExcludeField": g.shouldExcludeField,
		"ToOpenApiType":      toOpenApiType,
		"MapToModelType":     mapToModelType,
		"MapToApiType":       mapToApiType,
		"IsSimpleType":       utils.IsSimpleType,
		"IsComplexType":      utils.IsComplexType,
		"IsNullable":         isNullable,
		"Not":                not,
		"Types":              getTypesNamespace,
		"WrapID":             wrapID,
		"ShouldCreateField":  shouldCreateField,
		"GetGormQueryType":   getGormQueryType,
		"ToPtr":              toPtr,
		// will be replaced per model
		"ConvertToModel":           func() string { return "" },
		"ConvertToApi":             func() string { return "" },
		"ConvertToModelFromCreate": func() string { return "" },
		"ConvertToFilter":          func() string { return "" },
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

	for name, path := range g.cfg.UserTemplates {
		utpl := t.New(name)
		bytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error loading user-provided template %q: %w", name, err)
		}
		txt := string(bytes)
		_, err = utpl.Parse(txt)
		if err != nil {
			return fmt.Errorf("error parsing user-provided template %q: %w", name, err)
		}
	}

	return nil
}

func runNpxCommand(command string, args ...string) (string, error) {
	cmd := exec.Command("npx", append([]string{command}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
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
		if utils.IsSimpleType(field.Type) {
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

var primaryKeyRegex *regexp.Regexp = regexp.MustCompile("gorm:\"(.*?)\"")

// Whether the field should be excluded from create and update operations
func (g *generator) shouldExcludeField(field entity.GormModelField) bool {
	if slices.Contains(g.cfg.ExcludeFields, field.Name) {
		return true
	}

	if field.Parent != nil {
		if parentNamedType, ok := field.Parent.GetType().(*types.Named); ok {
			if parentNamedType.Obj().Pkg().Name() == "gorm" && parentNamedType.Obj().Name() == "Model" {
				return true
			}
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

func toOpenApiType(field entity.GormModelField) *utils.OpenApiType {
	if field.Tag != "" {
		settings, err := utils.ParseGoalesceTagSettings(field.Tag)
		if err == nil && len(settings) > 0 {
			openApiType := &utils.OpenApiType{}

			if typ, ok := settings["openapi_type"]; ok {
				openApiType.Type = typ
			}
			if ref, ok := settings["openapi_ref"]; ok {
				openApiType.Ref = &ref
			}
			if format, ok := settings["openapi_format"]; ok {
				openApiType.Format = &format
			}
			if nullable, ok := settings["openapi_nullable"]; ok {
				openApiType.Nullable = strings.ToLower(nullable) == "true"
			}

			return openApiType
		}
	}

	return utils.ToOpenApiType(field.Type)
}

func wrapID(model *entity.GormModelMetadata) string {
	result := "id"

	idField, err := utils.First(model.AllFields(), func(f *entity.GormModelField) bool {
		return f.Name == "ID"
	})
	if err == nil {
		if basicType, ok := idField.GetType().(*types.Basic); ok {
			wrapper := basicType.Name()
			result = fmt.Sprintf("%v(id)", wrapper)
		}
	}

	return result
}

func getGormQueryType(field *entity.GormModelField) string {
	t := field.GetGoType()
	if t == "time.Duration" {
		t = "int64"
	}
	return t
}

func toPtr(t string) string {
	if !strings.HasPrefix(t, "*") {
		t = "*" + t
	}
	return t
}

func not(v bool) bool {
	return !v
}

func isNullable(t string) bool {
	openApiType := utils.ToOpenApiType(t)
	return openApiType.Nullable
}

func shouldCreateField(f entity.GormModelField) bool {
	return !strings.HasPrefix(f.Type, "*") && !strings.HasPrefix(f.Type, "[]")

	// typ := f.GetType()

	// if f.MapFunc != nil {
	// 	mapFuncSig := f.GetMapFuncSignature()
	// 	resType := mapFuncSig.Results().At(0)
	// 	typ = resType.Type()
	// }

	// switch t := typ.(type) {
	// case *types.Basic:
	// 	return true
	// case *types.Pointer:
	// 	switch t.Elem().(type) {
	// 	case *types.Basic:
	// 		return true
	// 	default:
	// 		return false
	// 	}
	// case *types.Array, *types.Slice, *types.Map, *types.Chan, *types.Struct, *types.Tuple, *types.Signature, *types.Named, *types.Interface:
	// 	return false
	// default:
	// 	panic(fmt.Sprintf("impossible: %T", t))
	// }
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

func filepathToPackage(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
