package parse

import (
	"encoding/json"
	"testing"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Basic(t *testing.T) {
	actual, err := Parse("../../examples/basic/main.go")
	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	expectedTag := "`gorm:\"column:name;\"`"
	expected := &entity.GormModelMetadata{
		Name: "Person",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
				Tag:  &expectedTag,
			},
		},
	}

	assertJsonEq(t, expected, &actual[0])
}

func TestParse_Cars(t *testing.T) {
	actual, err := Parse("../../examples/cars/main.go")
	require.NoError(t, err)
	assert.Equal(t, 5, len(actual))

	expected := &entity.GormModelMetadata{
		Name: "Manufacturer",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Vehicles",
				Type: "[]Model",
			},
		},
	}
	assertJsonEq(t, expected, &actual[0])
	// TODO(joeriddles) assert all models in cars/main.go...
}

func assertJsonEq(t *testing.T, expected any, actual any) {
	actualBytes, err := json.Marshal(actual)
	require.NoError(t, err)
	actualJson := string(actualBytes)

	expectedBytes, err := json.Marshal(expected)
	require.NoError(t, err)
	expectedJson := string(expectedBytes)

	assert.JSONEq(t, expectedJson, actualJson)
}
