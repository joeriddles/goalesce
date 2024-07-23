package parse

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/joeriddles/gorm-oapi-codegen/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Basic(t *testing.T) {
	parser := NewParser(log.Default(), false)
	actual, err := parser.Parse("../../examples/basic/main.go")
	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	expectedIdTag := "`gorm:\"primarykey\"`"
	expectedDeletedAtTag := "`gorm:\"index\"`"
	expectedNameTag := "`gorm:\"column:name;\"`"

	expected := &entity.GormModelMetadata{
		Name:                "User",
		IsGormModelEmbedded: true,
		Fields: []*entity.GormModelField{
			{
				Name: "Model.ID",
				Type: "uint",
				Tag:  &expectedIdTag,
			},
			{
				Name: "Model.CreatedAt",
				Type: "time.Time",
			},
			{
				Name: "Model.UpdatedAt",
				Type: "time.Time",
			},
			{
				Name: "Model.DeletedAt",
				Type: "gorm.DeletedAt",
				Tag:  &expectedDeletedAtTag,
			},
			{
				Name: "Name",
				Type: "string",
				Tag:  &expectedNameTag,
			},
		},
	}

	assertJsonEq(t, expected, &actual[0])
}

func TestParse_Cars(t *testing.T) {
	parser := NewParser(log.Default(), false)
	actual, err := parser.Parse("../../examples/cars/main.go")
	require.NoError(t, err)
	assert.Equal(t, 5, len(actual))

	expectedIdTag := "`gorm:\"primarykey\"`"
	expectedDeletedAtTag := "`gorm:\"index\"`"

	expectedManufacturer := &entity.GormModelMetadata{
		Name:                "Manufacturer",
		IsGormModelEmbedded: true,
		Fields: []*entity.GormModelField{
			{
				Name: "Model.ID",
				Type: "uint",
				Tag:  &expectedIdTag,
			},
			{
				Name: "Model.CreatedAt",
				Type: "time.Time",
			},
			{
				Name: "Model.UpdatedAt",
				Type: "time.Time",
			},
			{
				Name: "Model.DeletedAt",
				Type: "gorm.DeletedAt",
				Tag:  &expectedDeletedAtTag,
			},
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
		Name:                "Model",
		IsGormModelEmbedded: true,
		Fields: []*entity.GormModelField{
			{
				Name: "Model.ID",
				Type: "uint",
				Tag:  &expectedIdTag,
			},
			{
				Name: "Model.CreatedAt",
				Type: "time.Time",
			},
			{
				Name: "Model.UpdatedAt",
				Type: "time.Time",
			},
			{
				Name: "Model.DeletedAt",
				Type: "gorm.DeletedAt",
				Tag:  &expectedDeletedAtTag,
			},
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

func TestParse_Custom(t *testing.T) {
	parser := NewParser(log.Default(), true)
	actual, err := parser.Parse("../../examples/custom/main.go")
	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	expectedIdTag := "`gorm:\"column:id;type:bigint;primaryKey;autoIncrement:true\" json:\"id\"`"
	expectedCreatedAtTag := "`gorm:\"column:created_at;type:timestamp with time zone\" json:\"created_at\"`"
	expectedUpdatedAtTag := "`gorm:\"column:updated_at;type:timestamp with time zone\" json:\"updated_at\"`"
	expectedDeletedAtTag := "`gorm:\"column:deleted_at;type:timestamp with time zone\" json:\"deleted_at\"`"

	expected := &entity.GormModelMetadata{
		Name:                "Custom",
		IsGormModelEmbedded: false,
		Fields: []*entity.GormModelField{
			{
				Name: "ID",
				Type: "int64",
				Tag:  &expectedIdTag,
			},
			{
				Name: "CreatedAt",
				Type: "time.Time",
				Tag:  &expectedCreatedAtTag,
			},
			{
				Name: "UpdatedAt",
				Type: "time.Time",
				Tag:  &expectedUpdatedAtTag,
			},
			{
				Name: "DeletedAt",
				Type: "gorm.DeletedAt",
				Tag:  &expectedDeletedAtTag,
			},
		},
	}

	assertJsonEq(t, expected, &actual[0])
}
