package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joeriddles/goalesce/examples/cars/generated/api"
	"github.com/joeriddles/goalesce/examples/cars/generated/repository"
	"github.com/joeriddles/goalesce/examples/cars/model"
	"github.com/joeriddles/goalesce/examples/cars/query"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newQuery(t *testing.T) *query.Query {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("PRAGMA foreign_keys = ON;") // enable FK constraints

	require.NoError(t, db.AutoMigrate(
		model.Manufacturer{},
		model.VehicleModel{},
		model.Vehicle{},
		model.VehicleForSale{},
		model.Part{},
		model.Person{},
	))
	query := query.Use(db)
	return query
}

func setupModels(t *testing.T, query *query.Query) (*model.Manufacturer, *model.VehicleModel, *model.Vehicle, *model.Person) {
	manufacturerRepo := repository.NewManufacturerRepository(query)
	manufacturer, err := manufacturerRepo.Create(context.Background(), model.Manufacturer{
		Name: "Mitsubishi",
	})
	require.NoError(t, err)

	vehicleModelRepo := repository.NewVehicleModelRepository(query)
	vehicleModel, err := vehicleModelRepo.Create(context.Background(), model.VehicleModel{
		ManufacturerID: manufacturer.ID,
		Name:           "Montero Sport",
	})
	require.NoError(t, err)

	personRepo := repository.NewPersonRepository(query)
	person, err := personRepo.Create(context.Background(), model.Person{
		Name: "Jim",
	})
	require.NoError(t, err)

	vehicleRepo := repository.NewVehicleRepository(query)
	vehicle, err := vehicleRepo.Create(context.Background(), model.Vehicle{
		VehicleModelID: vehicleModel.ID,
		PersonID:       person.ID,
		Vin:            "123",
	})
	require.NoError(t, err)

	return manufacturer, vehicleModel, vehicle, person
}

func Test_PostVehicle(t *testing.T) {
	// Arrange
	ctx := context.Background()
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	controller := api.NewVehicleController(query)

	_, vehicleModel, _, person := setupModels(t, query)
	vehicleModelID := int(vehicleModel.ID)
	personID := int(person.ID)

	// Act
	response, err := controller.PostVehicle(ctx, api.PostVehicleRequestObject{
		Body: &api.CreateVehicle{
			VehicleModelID: vehicleModelID,
			PersonID:       personID,
			Vin:            "456",
		},
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitPostVehicleResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 201, rec.Code)

	vehicle := &api.Vehicle{}
	err = json.Unmarshal(rec.Body.Bytes(), vehicle)
	require.NoError(t, err)
	assert.NotEqual(t, uint(0), vehicle.ID)
	assert.Equal(t, "456", vehicle.Vin)
	assert.Equal(t, vehicleModelID, vehicle.VehicleModelID)
	assert.Equal(t, personID, vehicle.PersonID)

	vehicleFromDb, err := repo.Get(ctx, int64(vehicle.ID))
	require.NoError(t, err)
	assert.Equal(t, "456", vehicleFromDb.Vin)
	assert.NotEqual(t, uint(0), vehicleFromDb.ID)
	assert.Equal(t, uint(vehicleModelID), vehicleFromDb.VehicleModelID)
	assert.Equal(t, uint(personID), vehicleFromDb.PersonID)
}

func Test_GetVehicle(t *testing.T) {
	// Arrange
	ctx := context.Background()
	query := newQuery(t)
	controller := api.NewVehicleController(query)

	_, _, vehicle, _ := setupModels(t, query)

	// Act
	response, err := controller.GetVehicle(ctx, api.GetVehicleRequestObject{})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitGetVehicleResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 200, rec.Code)

	vehicles := &[]api.Vehicle{}
	err = json.Unmarshal(rec.Body.Bytes(), vehicles)
	require.NoError(t, err)
	assert.Equal(t, 1, len(*vehicles))

	actualVehicle := &(*vehicles)[0]
	assert.Equal(t, vehicle.Vin, actualVehicle.Vin)
	assert.NotEqual(t, 0, actualVehicle.ID)
}

func Test_GetVehicle_WithFilters(t *testing.T) {
	// Arrange
	ctx := context.Background()
	query := newQuery(t)
	controller := api.NewVehicleController(query)

	_, _, vehicle, _ := setupModels(t, query)

	vehicleModelId := int(vehicle.VehicleModelID)
	personId := int(vehicle.PersonID)

	// Act
	testCases := []struct {
		name     string
		params   api.GetVehicleParams
		expected int
	}{
		{"Vin", api.GetVehicleParams{Vin: &vehicle.Vin}, 1},
		{"VehicleModelID", api.GetVehicleParams{VehicleModelID: &vehicleModelId}, 1},
		{"PersonID", api.GetVehicleParams{PersonID: &personId}, 1},
		{"Wrong Vin", api.GetVehicleParams{Vin: ptr("nope")}, 0},
		{"Wrong VehicleModelId", api.GetVehicleParams{VehicleModelID: ptr(vehicleModelId + 1)}, 0},
		{"Wrong PersonID", api.GetVehicleParams{PersonID: ptr(personId + 1)}, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := controller.GetVehicle(ctx, api.GetVehicleRequestObject{
				Params: tc.params,
			})
			require.NoError(t, err)

			// Assert
			rec := httptest.NewRecorder()
			err = response.VisitGetVehicleResponse(rec)
			require.NoError(t, err)
			assert.Equal(t, 200, rec.Code)

			vehicles := &[]api.Vehicle{}
			err = json.Unmarshal(rec.Body.Bytes(), vehicles)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, len(*vehicles))
		})
	}
}

func Test_GetVehicleID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	controller := api.NewVehicleController(query)

	_, _, vehicle, _ := setupModels(t, query)

	// Act
	response, err := controller.GetVehicleID(ctx, api.GetVehicleIDRequestObject{
		ID: int64(vehicle.ID),
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitGetVehicleIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 200, rec.Code)

	actual := &api.Vehicle{}
	err = json.Unmarshal(rec.Body.Bytes(), actual)
	require.NoError(t, err)
	assert.Equal(t, vehicle.Vin, actual.Vin)

	vehicleFromDb, err := repo.Get(ctx, int64(vehicle.ID))
	require.NoError(t, err)
	assert.Equal(t, vehicle.Vin, vehicleFromDb.Vin)
}

func Test_PutVehicleID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	controller := api.NewVehicleController(query)

	_, vehicleModel, vehicle, _ := setupModels(t, query)

	// Act
	response, err := controller.PutVehicleID(ctx, api.PutVehicleIDRequestObject{
		ID: int64(vehicle.ID),
		Body: &api.UpdateVehicle{
			VehicleModelID: int(vehicle.VehicleModelID),
			Vin:            "456",
		},
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitPutVehicleIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 204, rec.Code)

	vehicle, err = repo.Get(ctx, int64(vehicle.ID))
	require.NoError(t, err)
	assert.Equal(t, vehicleModel.ID, vehicle.VehicleModelID)
	assert.Equal(t, "456", vehicle.Vin)
}

func Test_DeleteVehicleID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	controller := api.NewVehicleController(query)

	_, _, vehicle, _ := setupModels(t, query)

	// Act
	response, err := controller.DeleteVehicleID(ctx, api.DeleteVehicleIDRequestObject{
		ID: int64(vehicle.ID),
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitDeleteVehicleIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 204, rec.Code)

	_, err = repo.Get(ctx, int64(vehicle.ID))
	require.Error(t, err, "Vehicle was not deleted from database")
}

func Test_PostVehicleForSale(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleForSaleRepository(query)
	controller := api.NewVehicleForSaleController(query)

	_, _, vehicle, _ := setupModels(t, query)
	vehicleID := int(vehicle.ID)

	// Act
	response, err := controller.PostVehicleForSale(context.Background(), api.PostVehicleForSaleRequestObject{
		Body: &api.CreateVehicleForSale{
			VehicleID: vehicleID,
			Amount:    "100.00",
			Duration:  60,
		},
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitPostVehicleForSaleResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 201, rec.Code)

	vehicleForSale := &api.VehicleForSale{}
	err = json.Unmarshal(rec.Body.Bytes(), vehicleForSale)
	require.NoError(t, err)
	assert.Equal(t, "100.00", vehicleForSale.Amount)
	assert.Equal(t, 60, vehicleForSale.Duration)
	assert.NotEqual(t, uint(0), vehicleForSale.ID)
	assert.Equal(t, vehicleID, vehicleForSale.VehicleID)

	vehicleForSaleFromDb, err := repo.Get(context.Background(), int64(vehicleForSale.ID))
	require.NoError(t, err)
	assert.NotEqual(t, uint(0), vehicleForSaleFromDb.ID)
	assert.Equal(t, vehicleID, int(vehicleForSaleFromDb.VehicleID))
	assert.Equal(t, time.Duration(60), vehicleForSaleFromDb.Duration)
	assert.True(t, vehicleForSaleFromDb.Amount.Equal(decimal.NewFromFloat(100.00)))
}

func ptr[T any](val T) *T {
	return &val
}
