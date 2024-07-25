package main

import (
	"testing"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/config"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/stretchr/testify/require"
)

// TODO(joeriddles): assert golden files

func Test_Basic(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/basic",
		OutputFile:      "../../examples/basic/generated",
		ModuleName:      "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:       "github.com/joeriddles/gorm-oapi-codegen/examples/basic",
		ClearOutputDir:  true,
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_Cars(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/cars",
		OutputFile:      "../../examples/cars/generated",
		ModuleName:      "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:       "github.com/joeriddles/gorm-oapi-codegen/examples/cars",
		ClearOutputDir:  true,
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_Custom(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath:   "../../examples/custom",
		OutputFile:        "../../examples/custom/generated",
		ModuleName:        "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:         "github.com/joeriddles/gorm-oapi-codegen/examples/custom",
		AllowCustomModels: true,
		ClearOutputDir:    true,
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_Circular(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath:   "../../examples/circular",
		OutputFile:        "../../examples/circular/generated",
		ModuleName:        "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:         "github.com/joeriddles/gorm-oapi-codegen/examples/circular",
		AllowCustomModels: true,
		ClearOutputDir:    true,
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_GenerateEcho(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/echo",
		OutputFile:      "../../examples/echo/generated",
		ModuleName:      "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:       "github.com/joeriddles/gorm-oapi-codegen/examples/echo",
		ClearOutputDir:  true,
		PruneYaml:       true,
		ModelsCodegen: &codegen.Configuration{
			PackageName: "api",
			Generate: codegen.GenerateOptions{
				Models: true,
			},
		},
		ServerCodegen: &codegen.Configuration{
			PackageName: "api",
			Generate: codegen.GenerateOptions{
				EchoServer:   true,
				Strict:       true,
				EmbeddedSpec: true,
			},
		},
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_GenerateExistingYaml(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/yaml",
		OutputFile:      "../../examples/yaml/generated",
		ModuleName:      "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:       "github.com/joeriddles/gorm-oapi-codegen/examples/yaml",
		ClearOutputDir:  true,
		PruneYaml:       true,
		OpenApiFile:     "../../examples/yaml/openapi.yaml",
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

// TODO(joeriddles): implement exclude
// func Test_GenerateExclude(t *testing.T) {
// 	cfg := &config.Config{
// 		InputFolderPath: "../../examples/exclude",
// 		OutputFile:      "../../examples/exclude/generated",
// 		ModuleName:      "github.com/joeriddles/gorm-oapi-codegen",
// 		ModelsPkg:       "github.com/joeriddles/gorm-oapi-codegen/examples/yaml",
// 		ClearOutputDir:  true,
// 		PruneYaml:       true,
// 		ExcludeModels: []string{
// 			"Manufacturer",
// 			"Model",
// 		},
// 	}
// 	require.NoError(t, cfg.Validate())
// 	err := run(cfg)
// 	require.NoError(t, err)
// }
