package custom

// Note: this does not inherit from gorm.Model
// See https://github.com/OAI/OpenAPI-Specification/issues/822
type Address struct {
	ID         int64 `gorm:"primaryKey;autoIncrement:true" json:"id"`
	City       string
	OccupantID int64
	Occupant   *Person
}

type Person struct {
	ID     int64 `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name   string
	HomeID int64
	Home   *Address
}
