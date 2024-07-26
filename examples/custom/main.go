package main

import (
	"github.com/joeriddles/goalesce/examples/custom/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "query",
		Mode:    gen.WithoutContext | gen.WithQueryInterface,
	})
	g.ApplyBasic(model.Custom{})
	g.Execute()
}
