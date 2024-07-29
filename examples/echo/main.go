package main

import (
	"github.com/joeriddles/goalesce/examples/echo/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "query",
		Mode:    gen.WithoutContext | gen.WithQueryInterface,
	})
	g.ApplyBasic(
		model.Manufacturer{},
		model.VehicleModel{},
		model.VehicleModel{},
		model.Part{},
		model.Person{},
	)
	g.Execute()
}
