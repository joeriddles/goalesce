package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// A vehicle manufacturer, like Chevrolet
type Manufacturer struct {
	gorm.Model
	Name     string
	Vehicles []VehicleModel
}

// A vehicle model, like a Chevrolet Silverado
type VehicleModel struct {
	gorm.Model
	Name           string
	ManufacturerID uint
	Manufacturer   Manufacturer
	Parts          []Part `gorm:"many2many:vehicle_parts;"`
}

// An individual of a model, like Joe's Chevrolet Silverado
type Vehicle struct {
	gorm.Model
	Vin            string
	VehicleModelID uint
	VehicleModel   VehicleModel
	PersonID       *int
	Person         *Person
}

// A vehicle for sale
type VehicleForSale struct {
	gorm.Model
	VehicleID uint
	Vehicle   Vehicle
	Amount    decimal.Decimal `goalesce:"openapi_type:string;map:MapVehicleForSaleAmount_Custom;map_api:MapApiVehicleForSaleAmount_Custom" gorm:"type:decimal(10,2);"`
	Duration  time.Duration   `goalesce:"openapi_type:integer;"`
}

// Custom map function for VehicleForSale.Amount from API to model
func MapVehicleForSaleAmount_Custom(apiAmount string) decimal.Decimal {
	result, err := decimal.NewFromString(apiAmount)
	if err != nil {
		panic(err)
	}
	return result
}

// Custom API map function for VehicleForSale.Amount from model to API
func MapApiVehicleForSaleAmount_Custom(amount decimal.Decimal) string {
	return amount.StringFixed(2)
}

// Custom map function for VehicleForSale.Duration from API to model
func MapVehicleForSaleDuration(apiDuration int) time.Duration {
	return time.Duration(int64(apiDuration))
}

// Custom API map function for VehicleForSale.Duration from model to API
func MapApiVehicleForSaleDuration(duration time.Duration) int {
	return int(int64(duration))
}

// A vehicle part for one or more models, like a muffler for all Chevrolet pickups
type Part struct {
	gorm.Model
	Name   string
	Cost   int
	Models []VehicleModel `gorm:"many2many:model_parts;"`
}

// A person, who may drive a vehicle
type Person struct {
	gorm.Model
	Name string
}
