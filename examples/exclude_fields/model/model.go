package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"column:name;"`
}

type Skill struct {
	gorm.Model
	Name   string `gorm:"column:name;"`
	UserID uint
	User   User
}
