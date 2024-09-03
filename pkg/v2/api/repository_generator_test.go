package api

import (
	"os"
	"testing"

	"github.com/joeriddles/goalesce/pkg/v2/config"
	"github.com/joeriddles/goalesce/pkg/v2/entity"
	"github.com/joeriddles/goalesce/pkg/v2/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryGenerator_Generate(t *testing.T) {
	// Arrange
	cfg, err := config.FromYamlFile("../../../examples/basic/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())

	model := &entity.GormModelMetadata{
		Name: "User",
		Fields: []*entity.GormModelField{
			{
				Name: "ID",
				Type: "int",
			},
			{
				Name: "Name",
				Type: "string",
				Tag:  "`gorm:\"column:name;\"`",
			},
		},
	}
	services := generator.NewGeneratorServices(cfg)
	generator := newRepositoryGenerator(model, services)

	// Act
	actual, err := generator.Generate()
	require.NoError(t, err, actual)

	if err == nil {
		filepath := "./test/repository_generator_test_generate.output.txt"
		require.NoError(t, os.WriteFile(filepath, []byte(actual), os.ModePerm))
		defer os.Remove(filepath)
	}

	// Assert
	expectedBytes, err := os.ReadFile("./test/repository_generator_test_generate.expected.txt")
	require.NoError(t, err)
	expected := string(expectedBytes)

	assert.Equal(t, expected, actual)

}
