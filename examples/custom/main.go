package custom

// Note: this does not inherit from gorm.Model
type Custom struct {
	Name string `gorm:"column:name;"`
}
