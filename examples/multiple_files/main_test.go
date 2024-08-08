package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/joeriddles/goalesce/examples/multiple_files/generated/api"
	"github.com/joeriddles/goalesce/examples/multiple_files/generated/repository"
	"github.com/joeriddles/goalesce/examples/multiple_files/model"
	"github.com/joeriddles/goalesce/examples/multiple_files/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newQuery(t *testing.T) *query.Query {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(model.User{}))
	query := query.Use(db)
	return query
}

func Test_PostUser(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewUserRepository(query)
	controller := api.NewUserController(query)

	// Act
	response, err := controller.PostUser(context.Background(), api.PostUserRequestObject{
		Body: &api.CreateUser{
			Name: "Bob",
		},
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitPostUserResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 201, rec.Code)

	user := &model.User{}
	err = json.Unmarshal(rec.Body.Bytes(), user)
	require.NoError(t, err)
	assert.Equal(t, "Bob", user.Name)
	require.NotEqual(t, uint(0), user.ID)

	userFromDb, err := repo.Get(context.Background(), int64(user.ID))
	require.NoError(t, err)
	assert.Equal(t, user.Name, userFromDb.Name)
}

func Test_GetUser(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewUserRepository(query)
	user, err := repo.Create(context.Background(), model.User{Name: "Bob"})
	require.NoError(t, err)

	controller := api.NewUserController(query)

	// Act
	response, err := controller.GetUser(context.Background(), api.GetUserRequestObject{})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitGetUserResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 200, rec.Code)

	users := &[]model.User{}
	err = json.Unmarshal(rec.Body.Bytes(), users)
	require.NoError(t, err)
	assert.Equal(t, 1, len(*users))

	actualUser := &(*users)[0]
	assert.Equal(t, user.Name, actualUser.Name)
	assert.NotEqual(t, 0, actualUser.ID)
}

func Test_GetUserID(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewUserRepository(query)
	user, err := repo.Create(context.Background(), model.User{Name: "Bob"})
	require.NoError(t, err)

	controller := api.NewUserController(query)

	// Act
	response, err := controller.GetUserID(context.Background(), api.GetUserIDRequestObject{
		ID: int64(user.ID),
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitGetUserIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 200, rec.Code)

	actual := &model.User{}
	err = json.Unmarshal(rec.Body.Bytes(), actual)
	require.NoError(t, err)
	assert.Equal(t, user.Name, actual.Name)

	userFromDb, err := repo.Get(context.Background(), int64(user.ID))
	require.NoError(t, err)
	assert.Equal(t, user.Name, userFromDb.Name)
}

func Test_PutUserID(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewUserRepository(query)
	user, err := repo.Create(context.Background(), model.User{Name: "Bob"})
	require.NoError(t, err)

	controller := api.NewUserController(query)

	// Act
	newName := "Jim"
	response, err := controller.PutUserID(context.Background(), api.PutUserIDRequestObject{
		ID: int64(user.ID),
		Body: &api.UpdateUser{
			Name: newName,
		},
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitPutUserIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 204, rec.Code)

	user, err = repo.Get(context.Background(), int64(user.ID))
	require.NoError(t, err)
	assert.Equal(t, newName, user.Name)
}

func Test_DeleteUserID(t *testing.T) {
	// Arrange
	query := newQuery(t)
	repo := repository.NewUserRepository(query)
	user, err := repo.Create(context.Background(), model.User{Name: "Bob"})
	require.NoError(t, err)

	controller := api.NewUserController(query)

	// Act
	response, err := controller.DeleteUserID(context.Background(), api.DeleteUserIDRequestObject{
		ID: int64(user.ID),
	})
	require.NoError(t, err)

	// Assert
	rec := httptest.NewRecorder()
	err = response.VisitDeleteUserIDResponse(rec)
	require.NoError(t, err)
	assert.Equal(t, 204, rec.Code)

	_, err = repo.Get(context.Background(), int64(user.ID))
	require.Error(t, err, "User was not deleted from database")
}
