package model

import (
	"gorm.io/gorm"
)

type Yaml struct {
	gorm.Model
	Name string
}