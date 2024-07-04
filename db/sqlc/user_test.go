package db

import (
	"context"
	"testing"
	"time"

	"github.com/rouclec/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	password := util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)

	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.FullName, arg.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	userFound, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)

	require.Equal(t, userFound.Email, user.Email)
	require.Equal(t, userFound.Username, user.Username)
	require.Equal(t, userFound.FullName, user.FullName)

	require.WithinDuration(t, userFound.CreatedAt, user.CreatedAt, time.Second)
	require.WithinDuration(t, userFound.PasswordChangedAt, user.PasswordChangedAt, time.Second)
}
