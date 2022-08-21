package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	assert.NotNil(t, store)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	logrus.WithFields(logrus.Fields{
		"acc1.Balance": acc1.Balance,
		"acc2.Balance": acc2.Balance,
	}).Info("Before transaction")

	n := 5
	amount := int64(10)

	errsCh := make(chan error)
	resultsCh := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("TX %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errsCh <- err
			resultsCh <- result
		}()
	}

	existed := make(map[int]bool)
	// check results
	for i := 0; i < n; i++ {
		err := <-errsCh
		assert.NoError(t, err)

		result := <-resultsCh
		assert.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		assert.NotEmpty(t, transfer)
		assert.Equal(t, acc1.ID, transfer.FromAccountID)
		assert.Equal(t, acc2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		assert.NoError(t, err)

		// check from entry
		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, acc1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		// check to entry
		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, acc2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, acc1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, acc2.ID, toAccount.ID)

		// check balances
		logrus.WithFields(logrus.Fields{
			"fromAccount.Balance": fromAccount.Balance,
			"toAccount.Balance":   toAccount.Balance,
		}).Info("Transaction")

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)

		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	assert.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	assert.NoError(t, err)

	logrus.WithFields(logrus.Fields{
		"updatedAcc1.Balance": updatedAcc1.Balance,
		"updatedAcc2.Balance": updatedAcc2.Balance,
	}).Info("After transaction")

	assert.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)
	assert.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	logrus.WithFields(logrus.Fields{
		"acc1.Balance": acc1.Balance,
		"acc2.Balance": acc2.Balance,
	}).Info("Before transaction")

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccID := acc1.ID
		toAccID := acc2.ID

		if i%2 == 1 {
			fromAccID = acc2.ID
			toAccID = acc1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccID,
				ToAccountID:   toAccID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	// check the final updated balance
	updatedAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	assert.NoError(t, err)

	updatedAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	assert.NoError(t, err)

	assert.Equal(t, acc1.Balance, updatedAcc1.Balance)
	assert.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
