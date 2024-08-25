package db

import (
	"context"
	"testing"
	"time"

	"github.com/kvgtl/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmailAddress(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), createdUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.HashedPassword, createdUser.HashedPassword)
	require.Equal(t, user.FullName, createdUser.FullName)
	require.Equal(t, user.Email, createdUser.Email)
	require.WithinDuration(t, user.PasswordChangedAt, createdUser.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user.CreatedAt, createdUser.CreatedAt, time.Second)
}
