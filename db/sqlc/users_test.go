package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/w0ikid/go-bank/util"
	"testing"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomString(10)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	err = util.CheckPassword(hashedPassword, password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.True(t, user1.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user1.CreatedAt)
}
