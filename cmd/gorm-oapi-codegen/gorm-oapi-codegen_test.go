package main

import (
	"testing"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/config"
	"github.com/stretchr/testify/require"
)

// TODO(joeriddles): assert golden files

func Test_Basic(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.WithInputFolderPath("../../examples/basic")
	require.NoError(t, err)
	err = cfg.WithOutputFile("./generated/basic")
	require.NoError(t, err)
	cfg.WithModuleName("github.com/joeriddles/gorm-oapi-codegen")
	cfg.WithModelPkg("github.com/joeriddles/gorm-oapi-codegen/examples/basic")
	cfg.WithClearOutputDir(true)

	err = run(cfg)
	require.NoError(t, err)
}

func Test_Cars(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.WithInputFolderPath("../../examples/cars")
	require.NoError(t, err)
	err = cfg.WithOutputFile("./generated/cars")
	require.NoError(t, err)
	cfg.WithModuleName("github.com/joeriddles/gorm-oapi-codegen")
	cfg.WithModelPkg("github.com/joeriddles/gorm-oapi-codegen/examples/cars")
	cfg.WithClearOutputDir(true)

	err = run(cfg)
	require.NoError(t, err)
}

func Test_Custom(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.WithInputFolderPath("../../examples/custom")
	require.NoError(t, err)
	err = cfg.WithOutputFile("./generated/custom")
	require.NoError(t, err)
	cfg.WithModuleName("github.com/joeriddles/gorm-oapi-codegen")
	cfg.WithModelPkg("github.com/joeriddles/gorm-oapi-codegen/examples/custom")
	cfg.WithAllowCustomModels(true)
	cfg.WithClearOutputDir(true)

	err = run(cfg)
	require.NoError(t, err)
}

func Test_Circular(t *testing.T) {
	cfg := config.NewConfig()
	err := cfg.WithInputFolderPath("../../examples/circular")
	require.NoError(t, err)
	err = cfg.WithOutputFile("./generated/circular")
	require.NoError(t, err)
	cfg.WithModuleName("github.com/joeriddles/gorm-oapi-codegen")
	cfg.WithModelPkg("github.com/joeriddles/gorm-oapi-codegen/examples/circular")
	cfg.WithAllowCustomModels(true)
	cfg.WithClearOutputDir(true)

	err = run(cfg)
	require.NoError(t, err)
}
