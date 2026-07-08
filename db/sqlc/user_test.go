package db

import (
	"context"
	"testing"
	"time"

	"github.com/NatdanaiKhe/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	username := util.RandomString(8)

	arg := CreateUserParams{
		Username:       username,
		HashedPassword: util.RandomString(12),
		FullName:       util.RandomString(12),
		Email:          util.RandomEmail(username),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	result, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.Username, result.Username)
	require.Equal(t, user.HashedPassword, result.HashedPassword)
	require.Equal(t, user.FullName, result.FullName)
	require.Equal(t, user.Email, result.Email)

	require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, result.PasswordChangedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserParams{
		Username:       user.Username,
		HashedPassword: util.RandomString(12),
		FullName:       util.RandomString(12),
		Email:          util.RandomEmail(user.Username),
	}
	result, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg.Username, result.Username)
	require.Equal(t, arg.HashedPassword, result.HashedPassword)
	require.Equal(t, arg.FullName, result.FullName)
	require.Equal(t, arg.Email, result.Email)

	require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
	require.NotEqual(t, user.PasswordChangedAt.IsZero(), result.PasswordChangedAt.IsZero())
}
