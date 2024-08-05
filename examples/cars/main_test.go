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

func Test_PostVehicle(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	controller := api.NewVehicleController(query)

	personRepo := repository.NewPersonRepository(query)
	person, err := personRepo.Create(context.Background(), model.Person{
		Name: "Jim",
	})
	require.NoError(t, err)
	personID := int(person.ID)

	// Act
	response, err := controller.PostVehicle(context.Background(), api.PostVehicleRequestObject{
		Body: &api.CreateVehicle{
			Vin:      "123",
			PersonID: &personID,
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
	assert.Equal(t, "123", vehicle.Vin)
	assert.NotEqual(t, uint(0), vehicle.ID)
	assert.Equal(t, personID, *vehicle.PersonID)

	vehicleFromDb, err := repo.Get(context.Background(), int64(vehicle.ID))
	require.NoError(t, err)
	assert.Equal(t, "123", vehicleFromDb.Vin)
	assert.NotEqual(t, uint(0), vehicleFromDb.ID)
	assert.Equal(t, personID, *vehicleFromDb.PersonID)
}

func Test_GetVehicle(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	vehicle, err := repo.Create(context.Background(), model.Vehicle{Vin: "123"})
	require.NoError(t, err)

	controller := api.NewVehicleController(query)

	// Act
	response, err := controller.GetVehicle(context.Background(), api.GetVehicleRequestObject{})
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

func Test_GetVehicleID(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	vehicle, err := repo.Create(context.Background(), model.Vehicle{Vin: "123"})
	require.NoError(t, err)

	controller := api.NewVehicleController(query)

	// Act
	response, err := controller.GetVehicleID(context.Background(), api.GetVehicleIDRequestObject{
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

	vehicleFromDb, err := repo.Get(context.Background(), int64(vehicle.ID))
	require.NoError(t, err)
	assert.Equal(t, vehicle.Vin, vehicleFromDb.Vin)
}

func Test_PutVehicleID(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	vehicle, err := repo.Create(context.Background(), model.Vehicle{Vin: "123"})
	require.NoError(t, err)

	controller := api.NewVehicleController(query)

	// Act
	newVin := "Jim"
	response, err := controller.PutVehicleID(context.Background(), api.PutVehicleIDRequestObject{
		ID: int64(vehicle.ID),
		Body: &api.UpdateVehicle{
			Vin: &newVin,
		},
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitPutVehicleIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 204, rec.Code)

	vehicle, err = repo.Get(context.Background(), int64(vehicle.ID))
	require.NoError(t, err)
	assert.Equal(t, newVin, vehicle.Vin)
}

func Test_DeleteVehicleID(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleRepository(query)
	vehicle, err := repo.Create(context.Background(), model.Vehicle{Vin: "123"})
	require.NoError(t, err)

	controller := api.NewVehicleController(query)

	// Act
	response, err := controller.DeleteVehicleID(context.Background(), api.DeleteVehicleIDRequestObject{
		ID: int64(vehicle.ID),
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitDeleteVehicleIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 204, rec.Code)

	_, err = repo.Get(context.Background(), int64(vehicle.ID))
	require.Error(t, err, "Vehicle was not deleted from database")
}

func Test_PostVehicleForSale(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewVehicleForSaleRepository(query)
	controller := api.NewVehicleForSaleController(query)

	vehicleRepo := repository.NewVehicleRepository(query)
	vehicle, err := vehicleRepo.Create(context.Background(), model.Vehicle{
		Vin: "123",
	})
	require.NoError(t, err)
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
