package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	assert.NotNil(t, store)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errsCh := make(chan error)
	resultsCh := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errsCh <- err
			resultsCh <- result
		}()
	}

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
	}

}
