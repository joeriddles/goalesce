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
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

// Generate from GORM model metadata
func Generate(templates embed.FS, metadatas []*entity.GormModelMetadata) error {
	t := template.New("gorm_oapi_codegen")
	if err := loadTemplates(templates, t); err != nil {
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
		Generate: codegen.GenerateOptions{
			StdHTTPServer: true,
			Strict:        true,
		},
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

func generateOpenApiRoutes(t *template.Template, metadata *entity.GormModelMetadata) error {
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

// loadTemplates loads all of our template files into a text/template. The
// path of template is relative to the templates directory.
func loadTemplates(src embed.FS, t *template.Template) error {
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

func runNpxCommand(command string, args ...string) (string, error) {
	cmd := exec.Command("npx", append([]string{command}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
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
