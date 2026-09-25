package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TxResultChan struct {
	res TransferTxResult
	err error
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	fmt.Println("Before Transaction >>>> ")
	fmt.Println("Acc1 balance >>> ", acc1.Balance)
	fmt.Println("Acc2 balance >>> ", acc2.Balance)

	// run n concurrent transfers
	n := 15
	amount := int64(10)

	// errs := make(chan error)
	// results := make(chan TransferTxResult)

	chres := make(chan TxResultChan, n)

	for range n {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId:   acc2.ID,
				Amount:        amount,
			})

			chres <- TxResultChan{res: result, err: err}
		}()
	}

	// check results
	existed := make(map[int]bool)

	for range n {
		out := <-chres
		require.NoError(t, out.err)

		res := out.res
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

		// check accounts
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)
		
		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc1.ID, toAccount.ID)

		// check accounts balances
		fmt.Println("at Transaction >>>> ")
		fmt.Println("from Acc balance >>> ", fromAccount.Balance)
		fmt.Println("to Acc balance >>> ", toAccount.Balance)

		// diff bw input acc1 balance and output fromacc balance
		diff1 := acc1.Balance - fromAccount.Balance

		// diff bw input toacc balance and output acc2 balance
		diff2 := toAccount.Balance - acc2.Balance

		// if transaction is correct then diff1 n diff2 shoudl be same
		require.Equal(t, diff1, diff2)

		// diff should be positive
		require.True(t, diff1 > 0)

		// diff should be dvisible by amt of money in each transaction i.e for 1st it should be decreased by 1*amt, for 2nd 2*amt and soc on
		require.True(t, diff1%amount==0)

		k:=int(diff1/amount)
		require.True(t, k>=1 && k<=n)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println("After Transaction >>>> ")
	fmt.Println("Acc1 balance >>> ", updatedAccount1.Balance)
	fmt.Println("Acc2 balance >>> ", updatedAccount2.Balance)

	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, acc1.Balance+int64(n)*amount, updatedAccount2.Balance)



}
