package main

import (
	"testing"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/stretchr/testify/require"
)

// TODO(joeriddles): assert golden files

func Test_Basic(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/basic",
		OutputFile:      "../../examples/basic/generated",
		ModuleName:      "github.com/joeriddles/goalesce/examples/basic",
		ModelsPkg:       "github.com/joeriddles/goalesce/examples/basic",
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
		ModuleName:      "github.com/joeriddles/goalesce/examples/cars",
		ModelsPkg:       "github.com/joeriddles/goalesce/examples/cars",
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
		ModuleName:        "github.com/joeriddles/goalesce/examples/custom",
		ModelsPkg:         "github.com/joeriddles/goalesce/examples/custom",
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
		ModuleName:        "github.com/joeriddles/goalesce/examples/circular",
		ModelsPkg:         "github.com/joeriddles/goalesce/examples/circular",
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
		ModuleName:      "github.com/joeriddles/goalesce/examples/echo",
		ModelsPkg:       "github.com/joeriddles/goalesce/examples/echo",
		ClearOutputDir:  true,
		PruneYaml:       true,
		TypesCodegen: &config.OApiGenConfiguration{
			Configuration: codegen.Configuration{
				PackageName: "api",
				Generate: codegen.GenerateOptions{
					Models: true,
				},
			},
		},
		ServerCodegen: &config.OApiGenConfiguration{
			Configuration: codegen.Configuration{
				PackageName: "api",
				Generate: codegen.GenerateOptions{
					EchoServer:   true,
					Strict:       true,
					EmbeddedSpec: true,
				},
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
		ModuleName:      "github.com/joeriddles/goalesce/examples/yaml",
		ModelsPkg:       "github.com/joeriddles/goalesce/examples/yaml",
		ClearOutputDir:  true,
		PruneYaml:       true,
		OpenApiFile:     "../../examples/yaml/openapi.yaml",
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_GenerateExclude(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/exclude",
		OutputFile:      "../../examples/exclude/generated",
		ModuleName:      "github.com/joeriddles/goalesce/examples/exclude",
		ModelsPkg:       "github.com/joeriddles/goalesce/examples/exclude",
		ClearOutputDir:  true,
		PruneYaml:       true,
		ExcludeModels: []string{
			"Manufacturer",
			"Model",
		},
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}

func Test_GenerateExclude_Yaml(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/exclude/config.yaml")
	require.NoError(t, err)

	err = run(cfg)
	require.NoError(t, err)
}

func Test_GenerateNestedTypes_Yaml(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/types/config.yaml")
	require.NoError(t, err)

	err = run(cfg)
	require.NoError(t, err)
}
