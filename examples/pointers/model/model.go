package model

type User struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"column:name;"`
}

type Pointer struct {
	ID int `gorm:"primaryKey"`

	UserID int
	User   User

	UserPtrID *int
	UserPtr   *User

	Users       []User   `gorm:"many2many:pointer_users;"`
	PtrUsers    *[]User  `gorm:"many2many:pointer_ptr_users;"`
	UserPtrs    []*User  `gorm:"many2many:pointer_user_ptrs;"`
	PtrUserPtrs *[]*User `gorm:"many2many:pointer_ptr_user_ptrs;"`
}
