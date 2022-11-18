package db

import (
	"context"
	"testing"
	"time"

	"github.com/harryng22/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(10),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	expectedUser := CreateRandomUser(t)
	actualUser, err := testQueries.GetUser(context.Background(), expectedUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, actualUser)
	require.Equal(t, expectedUser.Username, actualUser.Username)
	require.Equal(t, expectedUser.FullName, actualUser.FullName)
	require.Equal(t, expectedUser.Email, actualUser.Email)
	require.Equal(t, expectedUser.HashedPassword, actualUser.HashedPassword)
	require.WithinDuration(t, expectedUser.PasswordChangedAt, actualUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, expectedUser.CreatedAt, actualUser.CreatedAt, time.Second)
}
