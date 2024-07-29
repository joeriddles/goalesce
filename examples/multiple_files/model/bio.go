package model

import "gorm.io/gorm"

type Bio struct {
	gorm.Model
	Description string
	UserID      *uint
	User        User
}
