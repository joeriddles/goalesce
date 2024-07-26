package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp with time zone" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at"`
}

type Custom struct {
	Base
	Name string
}
