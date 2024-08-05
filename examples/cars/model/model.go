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
	Amount    decimal.Decimal `goalesce:"openapi_type:string" gorm:"type:decimal(10,2);"`
	Duration  time.Duration   `goalesce:"openapi_type:integer;"`
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
