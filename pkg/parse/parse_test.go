package parse

import (
	"encoding/json"
	"testing"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Basic(t *testing.T) {
	parser := NewParser()
	actual, err := parser.Parse("../../examples/basic/main.go")
	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	expectedTag := "`gorm:\"column:name;\"`"
	expected := &entity.GormModelMetadata{
		Name: "User",
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
	parser := NewParser()
	actual, err := parser.Parse("../../examples/cars/main.go")
	require.NoError(t, err)
	assert.Equal(t, 5, len(actual))

	expectedManufacturer := &entity.GormModelMetadata{
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
	expectedModelPartsTag := "`gorm:\"many2many:vehicle_parts;\"`"
	expectedModel := &entity.GormModelMetadata{
		Name: "Model",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "ManufacturerID",
				Type: "uint",
			},
			{
				Name: "Manufacturer",
				Type: "*Manufacturer",
			},
			{
				Name: "Parts",
				Type: "[]*Part",
				Tag:  &expectedModelPartsTag,
			},
		},
	}
	assertJsonEq(t, expectedManufacturer, &actual[0])
	assertJsonEq(t, expectedModel, &actual[1])
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
