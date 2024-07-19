package basic

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name string `gorm:"column:name;"`
}
