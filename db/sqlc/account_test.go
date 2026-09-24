package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dhirendraj-cmd/gosimplebankaccount/util"
	"github.com/stretchr/testify/require"
)


func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: util.RandOwner(),
		Balance: util.RandMoney(),
		Currency: util.RandCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}


func TestCreateAccount(t *testing.T){
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := CreateRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T){
	acc1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID: acc1.ID,
		Balance: util.RandMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T){
	acc1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)


	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}


func TestListAccounts(t *testing.T){
	for i:=0; i<151; i++{
		CreateRandomAccount(t)
	}

	arg:=ListAccountsParams{
		Limit: 9,
		Offset: 9,
	}

	acc, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, acc, 9)

	for _, account := range acc{
		require.NotEmpty(t, account)
	}

}
