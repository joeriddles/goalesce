package api

import (
	"os"
	"testing"

	"github.com/joeriddles/goalesce/pkg/v2/builder"
	"github.com/joeriddles/goalesce/pkg/v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainGenerator_Generate(t *testing.T) {
	// Arrange
	cfg, err := config.FromYamlFile("../../../examples/basic/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())

	goCodeBuilder := builder.NewGoCodeBuilder()
	generator, err := newMainGenerator(
		cfg,
		goCodeBuilder,
	)
	require.NoError(t, err)

	// Act
	actual, err := generator.Generate(nil)
	require.NoError(t, err, actual)

	if err == nil {
		filepath := "./test/main_generator_test_generate.output.txt"
		require.NoError(t, os.WriteFile(filepath, []byte(actual), os.ModePerm))
		defer os.Remove(filepath)
	}

	// Assert
	expectedBytes, err := os.ReadFile("./test/main_generator_test_generate.expected.txt")
	require.NoError(t, err)
	expected := string(expectedBytes)

	assert.Equal(t, expected, actual)
}
