package api

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/joeriddles/goalesce/pkg/v2/config"
	"github.com/joeriddles/goalesce/pkg/v2/entity"
	"github.com/joeriddles/goalesce/pkg/v2/logger"
	"github.com/joeriddles/goalesce/pkg/v2/parser"
	"github.com/joeriddles/goalesce/pkg/v2/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	codegen_util "github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

type ApiGenerator interface {
	Generate(metadatas []*entity.GormModelMetadata) error
}

type controllerGeneratorFactory = func(createApi *entity.GormModelMetadata, updateApi *entity.GormModelMetadata) controllerGenerator

type apiGenerator struct {
	cfg *config.Config

	logger logger.Logger

	parser parser.Parser

	openapiYamlControllerGenerator *openapiYamlControllerGenerator
	repositoryGenerator            *repositoryGenerator
	controllerGeneratorFactory     controllerGeneratorFactory
	serverGenerator                *serverGenerator
	mainGenerator                  *mainGenerator
}

func NewGenerator() (ApiGenerator, error) {
	return &apiGenerator{}, nil
}

func (g *apiGenerator) Generate(metadatas []*entity.GormModelMetadata) error {
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
	apiMetadatas, err := g.parser.Parse(filepath.Dir(g.cfg.TypesCodegen.OutputFile))
	if err != nil {
		return err
	}

	for _, metadata := range metadatas {
		apiMetadata, err := utils.First(apiMetadatas, func(m *entity.GormModelMetadata) bool {
			return m.Name == metadata.Name
		})
		if err != nil {
			g.logger.Log("could not find apiMetadata for %v", metadata.Name)
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
		if err := g.generateController(metadata, createApiMetadata, updateApiMetadata); err != nil {
			return err
		}
	}

	filteredMetadatas := []*entity.GormModelMetadata{}
	for _, metadata := range metadatas {
		if !slices.Contains(g.cfg.ExcludeModels, metadata.Name) {
			filteredMetadatas = append(filteredMetadatas, metadata)
		}
	}
	if err := g.generateServer(filteredMetadatas); err != nil {
		return err
	}

	if g.cfg.GenerateMain {
		if err := g.generateMain(); err != nil {
			return err
		}
	}

	// if err := g.generateMapperUtil(); err != nil {
	// 	return err
	// }

	return nil
}

func (g *apiGenerator) generateOpenApiYaml(metadatas []*entity.GormModelMetadata) error {
	for _, metadata := range metadatas {
		_, err := g.generateOpenApiRoutes(metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *apiGenerator) runCodegenTool() error {
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

func (g *apiGenerator) generateOpenApiRoutes(metadata *entity.GormModelMetadata) (string, error) {
	fp := filepath.Join(g.cfg.OutputFile, fmt.Sprintf("%v.gen.yaml", utils.ToSnakeCase(metadata.Name)))
	code, err := g.openapiYamlControllerGenerator.Generate(metadata)
	if err != nil {
		return "", err
	}
	return fp, g.generateFile(fp, code)
}

func (g *apiGenerator) generateMapper(
	metadata *entity.GormModelMetadata,
	apiMetadata *entity.GormModelMetadata,
) error {
	return nil
}

func (g *apiGenerator) generateRepository(metadata *entity.GormModelMetadata) error {
	filename := fmt.Sprintf("%v_repository.gen.go", utils.ToSnakeCase(metadata.Name))
	fp := filepath.Join(g.cfg.RepositoryConfiguration.OutputFile, filename)
	code, err := g.repositoryGenerator.Generate(metadata)
	if err != nil {
		return err
	}
	return g.generateFile(fp, code)
}

func (g *apiGenerator) generateController(
	metadata *entity.GormModelMetadata,
	createApiMetadata *entity.GormModelMetadata,
	updateApiMetadata *entity.GormModelMetadata,
) error {
	fp := filepath.Join(g.cfg.OutputFile, "api", fmt.Sprintf("%v_controller.gen.go", utils.ToSnakeCase(metadata.Name)))
	controllerGenerator := g.controllerGeneratorFactory(createApiMetadata, updateApiMetadata)
	code, err := controllerGenerator.Generate(metadata)
	if err != nil {
		return err
	}
	return g.generateFile(fp, code)
}

func (g *apiGenerator) generateServer(metadatas []*entity.GormModelMetadata) error {
	_, err := g.serverGenerator.Generate(metadatas)
	return err
}

func (g *apiGenerator) generateMain() error {
	_, err := g.mainGenerator.Generate(nil) // TODO: refactor Generator interface
	return err
}

func (g *apiGenerator) generateFile(filepath, content string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(content)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
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
