package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/joeriddles/goalesce/examples/basic/generated/api"
	"github.com/joeriddles/goalesce/examples/basic/generated/repository"
	"github.com/joeriddles/goalesce/examples/basic/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Post(t *testing.T) {
	// Arrange
	controller := api.NewUserController()

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
}

func Test_List(t *testing.T) {
	// Arrange
	repo := repository.NewUserRepository()
	repo.Create(context.Background(), model.User{Name: "Bob"})

	controller := api.NewUserController()

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
	assert.Equal(t, "Bob", (*users)[0].Name)
}
