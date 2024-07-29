package parse

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/joeriddles/goalesce/pkg/config"
	"github.com/joeriddles/goalesce/pkg/entity"
	"github.com/joeriddles/goalesce/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopWriter struct{}

func (w *noopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

var noopLogger *log.Logger = log.New(&noopWriter{}, "", 0)

func TestParse_Basic(t *testing.T) {
	cfg := &config.Config{
		AllowCustomModels: false,
		InputFolderPath:   "../../examples/basic",
	}
	parser := NewParser(noopLogger, cfg)
	actual, err := parser.Parse("../../examples/basic/model")
	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	expected := &entity.GormModelMetadata{
		Name: "User",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
				Tag:  "gorm:\"column:name;\"",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			{
				Name:     "",
				Embedded: []*entity.GormModelMetadata{},
				Fields: []*entity.GormModelField{
					{
						Name: "ID",
						Type: "uint",
						Tag:  "gorm:\"primarykey\"",
					},
					{
						Name: "CreatedAt",
						Type: "time.Time",
					},
					{
						Name: "UpdatedAt",
						Type: "time.Time",
					},
					{
						Name: "DeletedAt",
						Type: "gorm.io/gorm.DeletedAt",
						Tag:  "gorm:\"index\"",
					},
				},
			},
		},
	}

	assertJsonEq(t, expected, &actual[0])
}

func TestParse_Cars(t *testing.T) {
	cfg := &config.Config{
		AllowCustomModels: false,
		InputFolderPath:   "../../examples/cars",
	}
	parser := NewParser(noopLogger, cfg)
	actual, err := parser.Parse("../../examples/cars/model")
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
				Type: "[]github.com/joeriddles/goalesce/examples/cars/model.VehicleModel",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			{
				Name:     "",
				Embedded: []*entity.GormModelMetadata{},
				Fields: []*entity.GormModelField{
					{
						Name: "ID",
						Type: "uint",
						Tag:  "gorm:\"primarykey\"",
					},
					{
						Name: "CreatedAt",
						Type: "time.Time",
					},
					{
						Name: "UpdatedAt",
						Type: "time.Time",
					},
					{
						Name: "DeletedAt",
						Type: "gorm.io/gorm.DeletedAt",
						Tag:  "gorm:\"index\"",
					},
				},
			},
		},
	}
	expectedModel := &entity.GormModelMetadata{
		Name: "VehicleModel",
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
				Type: "github.com/joeriddles/goalesce/examples/cars/model.Manufacturer",
			},
			{
				Name: "Parts",
				Type: "[]github.com/joeriddles/goalesce/examples/cars/model.Part",
				Tag:  "gorm:\"many2many:vehicle_parts;\"",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			{
				Name:     "",
				Embedded: []*entity.GormModelMetadata{},
				Fields: []*entity.GormModelField{
					{
						Name: "ID",
						Type: "uint",
						Tag:  "gorm:\"primarykey\"",
					},
					{
						Name: "CreatedAt",
						Type: "time.Time",
					},
					{
						Name: "UpdatedAt",
						Type: "time.Time",
					},
					{
						Name: "DeletedAt",
						Type: "gorm.io/gorm.DeletedAt",
						Tag:  "gorm:\"index\"",
					},
				},
			},
		},
	}

	actualManufacturer, err := utils.First(actual, func(f *entity.GormModelMetadata) bool {
		return f.Name == "Manufacturer"
	})
	require.NoError(t, err)
	assertJsonEq(t, expectedManufacturer, actualManufacturer)

	actualModel, err := utils.First(actual, func(f *entity.GormModelMetadata) bool {
		return f.Name == "VehicleModel"
	})
	require.NoError(t, err)
	assertJsonEq(t, expectedModel, actualModel)
	// TODO(joeriddles) assert all models in cars/main.go...
}

func TestParse_Custom(t *testing.T) {
	cfg := &config.Config{
		AllowCustomModels: true,
		InputFolderPath:   "../../examples/custom",
	}
	parser := NewParser(noopLogger, cfg)
	actual, err := parser.Parse("../../examples/custom/model")
	require.NoError(t, err)
	assert.Equal(t, 2, len(actual))

	expected := &entity.GormModelMetadata{
		Name:     "Base",
		Embedded: []*entity.GormModelMetadata{},
		Fields: []*entity.GormModelField{
			{
				Name: "ID",
				Type: "int64",
				Tag:  "gorm:\"column:id;type:bigint;primaryKey;autoIncrement:true\" json:\"id\"",
			},
			{
				Name: "CreatedAt",
				Type: "time.Time",
				Tag:  "gorm:\"column:created_at;type:timestamp with time zone\" json:\"created_at\"",
			},
			{
				Name: "UpdatedAt",
				Type: "time.Time",
				Tag:  "gorm:\"column:updated_at;type:timestamp with time zone\" json:\"updated_at\"",
			},
			{
				Name: "DeletedAt",
				Type: "gorm.io/gorm.DeletedAt",
				Tag:  "gorm:\"column:deleted_at;type:timestamp with time zone\" json:\"deleted_at\"",
			},
		},
	}
	assertJsonEq(t, expected, &actual[0])

	expected = &entity.GormModelMetadata{
		Name: "Custom",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			{
				Name:     "",
				Embedded: []*entity.GormModelMetadata{},
				Fields: []*entity.GormModelField{
					{
						Name: "ID",
						Type: "int64",
						Tag:  "gorm:\"column:id;type:bigint;primaryKey;autoIncrement:true\" json:\"id\"",
					},
					{
						Name: "CreatedAt",
						Type: "time.Time",
						Tag:  "gorm:\"column:created_at;type:timestamp with time zone\" json:\"created_at\"",
					},
					{
						Name: "UpdatedAt",
						Type: "time.Time",
						Tag:  "gorm:\"column:updated_at;type:timestamp with time zone\" json:\"updated_at\"",
					},
					{
						Name: "DeletedAt",
						Type: "gorm.io/gorm.DeletedAt",
						Tag:  "gorm:\"column:deleted_at;type:timestamp with time zone\" json:\"deleted_at\"",
					},
				},
			},
		},
	}
	assertJsonEq(t, expected, &actual[1])
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
