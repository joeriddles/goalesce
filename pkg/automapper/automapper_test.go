package automapper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAutoMapper_int(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(0, 0)
	from := 1
	to := -1

	// Assert
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, to)
}

type User struct {
	Name string
}

type User2 struct {
	Name string
}

func TestAutoMapper_Struct(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(User{}, User2{})
	from := User{Name: "Bob"}
	to := User2{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := User2{Name: "Bob"}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

type Cart struct {
	Items []string
}

type Cart2 struct {
	Items []string
}

func TestAutoMapper_StructWithList(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(Cart{}, Cart2{})
	from := Cart{Items: []string{"apple", "orange"}}
	to := Cart2{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := Cart2{Items: []string{"apple", "orange"}}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

type DifferentCase struct {
	UserId int
}

type DifferentCase2 struct {
	UserID int
}

func TestAutoMapper_StructDifferentCase(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(DifferentCase{}, DifferentCase2{})
	from := DifferentCase{UserId: 1}
	to := DifferentCase2{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := DifferentCase2{UserID: 1}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

type DifferentType struct {
	UserId int
}

type DifferentType2 struct {
	UserId int64
}

func TestAutoMapper_StructDifferentType(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(DifferentType{}, DifferentType2{})
	from := DifferentType{UserId: 1}
	to := DifferentType2{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := DifferentType2{UserId: 1}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

type NestedType struct {
	User User
}

type NestedType2 struct {
	User User2
}

func TestAutoMapper_StructNestedType(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(NestedType{}, NestedType2{})
	from := NestedType{User: User{Name: "Bob"}}
	to := NestedType2{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := NestedType2{User: User2{Name: "Bob"}}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

type NestedTypePointer struct {
	User *User
}

type NestedTypePointer2 struct {
	User *User2
}

func TestAutoMapper_StructNestedPointerType(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(NestedTypePointer{}, NestedTypePointer2{})
	from := NestedTypePointer{User: &User{Name: "Bob"}}
	to := NestedTypePointer2{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := NestedTypePointer2{User: &User2{Name: "Bob"}}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

type CreateUser struct {
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ModelUser struct {
	gorm.Model
	Name string `gorm:"column:name;"`
}

func TestAutoMapper_CreateUserToModelUser(t *testing.T) {
	// Arrange
	mapper := NewAutoMapper(CreateUser{}, ModelUser{})
	from := CreateUser{Name: "Bob"}
	to := ModelUser{}

	// Act
	err := mapper.MapTo(&from, &to)

	// Assert
	require.NoError(t, err)

	expected := ModelUser{Name: "Bob"}
	assert.Equal(t, expected, to)
	assertJsonEq(t, expected, to)
}

func assertJsonEq(t *testing.T, expected any, actual any) {
	actualBytes, err := json.Marshal(actual)
	require.NoError(t, err)
	actualJson := string(actualBytes)

	expectedBytes, err := json.Marshal(expected)
	require.NoError(t, err)
	expectedJson := string(expectedBytes)

	assert.JSONEq(t, expectedJson, actualJson)
}
