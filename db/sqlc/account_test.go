package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/NguyenMinhKhanhBK/simple_bank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)
	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, acc2)

	assert.Equal(t, acc1, acc2)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, acc2)

	assert.Equal(t, acc1.ID, acc2.ID)
	assert.Equal(t, acc1.Owner, acc2.Owner)
	assert.Equal(t, arg.Balance, acc2.Balance)
	assert.Equal(t, acc1.Currency, acc2.Currency)
	assert.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, 1*time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	assert.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, acc2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)

	for _, acc := range accounts {
		assert.NotEmpty(t, acc)
	}
}
