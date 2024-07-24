package main

import (
	"testing"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/config"
	"github.com/stretchr/testify/require"
)

// TODO(joeriddles): assert golden files

func Test_Basic(t *testing.T) {
	cfg := &config.Config{
		InputFolderPath: "../../examples/basic",
		OutputFile:      "./generated/basic",
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
		OutputFile:      "./generated/cars",
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
		OutputFile:        "./generated/custom",
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
		OutputFile:        "./generated/circular",
		ModuleName:        "github.com/joeriddles/gorm-oapi-codegen",
		ModelsPkg:         "github.com/joeriddles/gorm-oapi-codegen/examples/circular",
		AllowCustomModels: true,
		ClearOutputDir:    true,
	}
	require.NoError(t, cfg.Validate())
	err := run(cfg)
	require.NoError(t, err)
}
