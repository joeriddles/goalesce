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
	cfg, err := config.FromYamlFile("../../examples/basic/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())

	parser := NewParser(noopLogger, cfg)
	actual, err := parser.Parse(cfg.InputFolderPath)
	require.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	expected := &entity.GormModelMetadata{
		Pkg:  "model",
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
				Pkg:      "gorm",
				Name:     "Model",
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
	cfg, err := config.FromYamlFile("../../examples/cars/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())

	parser := NewParser(noopLogger, cfg)
	actual, err := parser.Parse(cfg.InputFolderPath)
	require.NoError(t, err)
	assert.Equal(t, 6, len(actual))

	expectedManufacturer := &entity.GormModelMetadata{
		Pkg:  "model",
		Name: "Manufacturer",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
			},
			{
				Name: "Vehicles",
				Type: "[]VehicleModel",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			{
				Pkg:      "gorm",
				Name:     "Model",
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
		Pkg:  "model",
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
				Type: "Manufacturer",
			},
			{
				Name: "Parts",
				Type: "[]Part",
				Tag:  "gorm:\"many2many:vehicle_parts;\"",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			{
				Pkg:      "gorm",
				Name:     "Model",
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
	cfg, err := config.FromYamlFile("../../examples/custom/config.yaml")
	require.NoError(t, err)
	require.NoError(t, cfg.Validate())

	parser := NewParser(noopLogger, cfg)
	actual, err := parser.Parse(cfg.InputFolderPath)
	require.NoError(t, err)
	assert.Equal(t, 2, len(actual))

	expectedBase := &entity.GormModelMetadata{
		Pkg:      "model",
		Name:     "Base",
		Embedded: []*entity.GormModelMetadata{},
		Fields: []*entity.GormModelField{
			{
				Name: "ID",
				Type: "int64",
				Tag:  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true"               json:"id"`,
			},
			{
				Name: "CreatedAt",
				Type: "time.Time",
				Tag:  `gorm:"column:created_at;type:timestamp with time zone;autoCreateTime;"   json:"created_at"`,
			},
			{
				Name: "UpdatedAt",
				Type: "time.Time",
				Tag:  `gorm:"column:updated_at;type:timestamp with time zone;isAutoUpdateTime;" json:"updated_at"`,
			},
			{
				Name: "DeletedAt",
				Type: "gorm.io/gorm.DeletedAt",
				Tag:  `gorm:"column:deleted_at;type:timestamp with time zone"                   json:"deleted_at"`,
			},
		},
	}
	assertJsonEq(t, expectedBase, &actual[0])

	expectedCustom := &entity.GormModelMetadata{
		Pkg:  "model",
		Name: "Custom",
		Fields: []*entity.GormModelField{
			{
				Name: "Name",
				Type: "string",
			},
		},
		Embedded: []*entity.GormModelMetadata{
			expectedBase,
		},
	}
	assertJsonEq(t, expectedCustom, &actual[1])
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
