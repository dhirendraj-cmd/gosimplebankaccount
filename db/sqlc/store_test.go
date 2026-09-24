package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)


type TxResultChan struct{
	res TransferTxResult
	err error
}


func TestTransferTx(t *testing.T){
	store := NewStore(testDB)

	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	// run n concurrent transfers
	n:=5
	amount := int64(10)

	// errs := make(chan error)
	// results := make(chan TransferTxResult)

	chres := make(chan TxResultChan, n)

	for range n{
		go func ()  {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId: acc2.ID,
				Amount: amount,
			})

			chres <- TxResultChan{res: result, err: err}
		}()
	}

	// check results
	for range n{
		out := <-chres
		require.NoError(t, out.err)

		res:=out.res
		require.NotEmpty(t, res)

		// check transfer
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)


		// check entries
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)


		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)


		// check account balances

	}

}

