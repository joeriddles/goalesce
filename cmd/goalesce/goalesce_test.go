package main

import (
	"testing"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/stretchr/testify/require"
)

// TODO(joeriddles): assert golden files

func Test_Basic(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/basic/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_Cars(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/cars/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_Custom(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/custom/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_Circular(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/circular/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_GenerateEcho(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/echo/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_GenerateExistingYaml(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/yaml/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_GenerateExclude(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/exclude/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_GenerateNestedTypes_Yaml(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/types/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_Generate_RepositoryConfig(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/repository/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}

func Test_Generate_MultipleModelFiles(t *testing.T) {
	cfg, err := config.FromYamlFile("../../examples/multiple_files/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())
	err = run(cfg)
	require.NoError(t, err)
}
