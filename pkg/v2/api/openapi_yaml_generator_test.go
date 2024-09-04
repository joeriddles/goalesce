package api

import (
	"os"
	"testing"

	"github.com/joeriddles/goalesce/pkg/v2/config"
	"github.com/joeriddles/goalesce/pkg/v2/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenapiYamlControllerGenerator_Generate(t *testing.T) {
	// Arrange
	cfg, err := config.FromYamlFile("../../../examples/basic/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())

	gormModel := &entity.GormModelMetadata{
		Pkg:      "gorm",
		Name:     "Model",
		Embedded: []*entity.GormModelMetadata{},
	}
	gormModel.Fields = []*entity.GormModelField{
		{
			Name:   "ID",
			Type:   "uint",
			Tag:    "gorm:\"primarykey\"",
			Parent: gormModel,
		},
		{
			Name:   "CreatedAt",
			Type:   "time.Time",
			Parent: gormModel,
		},
		{
			Name:   "UpdatedAt",
			Type:   "time.Time",
			Parent: gormModel,
		},
		{
			Name:   "DeletedAt",
			Type:   "gorm.io/gorm.DeletedAt",
			Tag:    "gorm:\"index\"",
			Parent: gormModel,
		},
	}

	model := &entity.GormModelMetadata{
		Pkg:  "model",
		Name: "User",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
				Tag:  "`gorm:\"column:name;\"`",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			gormModel,
		},
	}

	generator := newOpenapiYamlControllerGenerator()

	// Act
	actual, err := generator.Generate(model)
	require.NoError(t, err, actual)

	if err == nil {
		filepath := "./test/openapi_yaml_controller_generator_test_generate.output.yaml"
		require.NoError(t, os.WriteFile(filepath, []byte(actual), os.ModePerm))
		// defer os.Remove(filepath)
	}

	// Assert
	expectedBytes, err := os.ReadFile("./test/openapi_yaml_controller_generator_test_generate.expected.yaml")
	require.NoError(t, err)
	expected := string(expectedBytes)

	assert.Equal(t, expected, actual)

}
