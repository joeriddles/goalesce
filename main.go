package main

import (
	"bufio"
	"embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

// Embed the templates directory
//
//go:embed templates
var templates embed.FS

func errExit(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		errExit("Please specify a path to a folder of GORM models\n")
	} else if flag.NArg() > 1 {
		errExit("Only one folder path is accepted and it must be the last CLI argument\n")
	}

	folderPath := flag.Arg(0)
	if err := CodeGen(folderPath); err != nil {
		errExit(err.Error())
	}
}

func CodeGen(folderPath string) error {
	folderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return err
	}

	// Check path exists and we have permission to read it
	if _, err := os.Stat(folderPath); err != nil {
		return err
	}

	entries, err := os.ReadDir(folderPath)
	for _, entry := range entries {
		filename := entry.Name()
		entryFilepath := filepath.Join(folderPath, filename)
		metadatas, err := parse(entryFilepath)
		if err != nil {
			return err
		}
		if err := generate(metadatas); err != nil {
			return err
		}
	}
	return err
}

// Generate from GORM model metadata
func generate(metadatas []*GormModelMetadata) error {
	t := template.New("gorm_oapi_codegen")
	if err := LoadTemplates(templates, t); err != nil {
		return err
	}

	if err := os.RemoveAll("./generated"); err != nil {
		return err
	}

	if err := os.Mkdir("./generated", os.ModePerm); err != nil {
		if err.Error() != "mkdir generated: file exists" {
			return err
		}
	}

	for _, metadata := range metadatas {
		if err := generateOpenApiRoutes(t, metadata); err != nil {
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
	if err = os.WriteFile("./generated/types.gen.go", []byte(code), 0o644); err != nil {
		return err
	}

	code, err = codegen.Generate(swagger, codegen.Configuration{
		PackageName: "api",
		Generate:    codegen.GenerateOptions{StdHTTPServer: true},
	})

	if err != nil {
		return err
	}
	if err = os.WriteFile("./generated/server.gen.go", []byte(code), 0o644); err != nil {
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

	if _, err = runNpxCommand(
		"@redocly/openapi-cli@latest",
		"bundle",
		"./generated/openapi_base.gen.yaml",
		"-o",
		"./generated/openapi.yaml",
	); err != nil {
		return err
	}

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

func runNpxCommand(command string, args ...string) (string, error) {
	cmd := exec.Command("npx", append([]string{command}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func generateOpenApiBase(t *template.Template, metadatas []*GormModelMetadata) error {
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

func generateOpenApiRoutes(t *template.Template, metadata *GormModelMetadata) error {
	f, err := os.Create(fmt.Sprintf("generated/%v.gen.yaml", strings.ToLower(metadata.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.ExecuteTemplate(w, "openapi_controller.yaml", &metadata)
	w.Flush()

	return nil
}

// LoadTemplates loads all of our template files into a text/template. The
// path of template is relative to the templates directory.
func LoadTemplates(src embed.FS, t *template.Template) error {
	funcMap := template.FuncMap{
		"ToLower":       strings.ToLower,
		"ToOpenApiType": toOpenApiType,
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

// Parse GORM model metadata from the Go file
func parse(filepath string) ([]*GormModelMetadata, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, nil, parser.SkipObjectResolution) // ParseComments
	if err != nil {
		return nil, err
	}

	metadatas := []*GormModelMetadata{}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			metadata, err := parseGormModel(x)
			if err != nil {
				panic(err)
			}
			metadatas = append(metadatas, metadata)
		}
		return true
	})

	return metadatas, nil
}

type GormModelMetadata struct {
	Name   string
	Fields []*GormModelField
}

type GormModelField struct {
	Name string
	Type string
	Tag  *string
}

// Parse metadata about the GORM model node
func parseGormModel(node *ast.TypeSpec) (*GormModelMetadata, error) {
	if !checkIsGormModel(node) {
		return nil, errors.New("not a GORM model")
	}

	name := node.Name.Name
	fields := parseGormModelFields(node)

	metadata := &GormModelMetadata{
		Name:   name,
		Fields: fields,
	}
	return metadata, nil
}

func parseGormModelFields(node *ast.TypeSpec) []*GormModelField {
	fields := []*GormModelField{}
	ast.Inspect(node, func(n ast.Node) bool {
		switch f := n.(type) {
		case *ast.Field:
			if len(f.Names) == 0 {
				break
			}

			fName := f.Names[0].Name // TODO(joeriddles): support multiple names

			var fType string
			if typeId, ok := f.Type.(*ast.Ident); ok {
				fType = typeId.Name
			}

			var fTag *string
			if f.Tag != nil {
				fTag = &f.Tag.Value
			}

			field := &GormModelField{
				Name: fName,
				Type: fType,
				Tag:  fTag,
			}
			fields = append(fields, field)
		}
		return true
	})
	return fields
}

// Check if the ast node is a GORM model
func checkIsGormModel(node *ast.TypeSpec) bool {
	var isGorm = false

	ast.Inspect(node, func(n ast.Node) bool {
		switch f := n.(type) {
		case *ast.Field:
			if expr, ok := f.Type.(*ast.SelectorExpr); ok {
				xId, xIdOk := expr.X.(*ast.Ident)
				if xIdOk && xId.Name == "gorm" && expr.Sel.Name == "Model" {
					isGorm = true
				}
			}
		}
		return true
	})

	return isGorm
}

func toOpenApiType(t string) string {
	switch t {
	case "string":
		return "string"
	case "uint":
		return "integer"
	default:
		return "object"
	}
}
