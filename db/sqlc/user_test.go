package db

import (
	"context"
	"testing"
	"time"

	"github.com/NguyenMinhKhanhBK/simple_bank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	assert.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)
	assert.NotZero(t, user.CreatedAt)
	assert.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
	assert.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
