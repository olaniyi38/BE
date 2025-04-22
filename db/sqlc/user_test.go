package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/olaniyi38/BE/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	randomUserName := util.RandomString(10)
	randomEmail := fmt.Sprintf("%v@gmail.com", util.RandomString(6))
	randomPassword, err := util.GeneratePassword(util.RandomString(10))
	randomFullName := fmt.Sprintf("%v %v", util.RandomString(8), util.RandomString(8))
	require.NoError(t, err)

	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username: randomUserName,
		Email:    randomEmail,
		Password: randomPassword,
		FullName: randomFullName,
	})
	require.NoError(t, err)

	require.NotEmpty(t, user)
	require.Equal(t, randomUserName, user.Username)
	require.Equal(t, randomEmail, user.Email)
	require.Equal(t, randomPassword, user.Password)
	require.NotZero(t, user.FullName, randomFullName)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	gotuser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)

	require.Equal(t, gotuser, user)
}
