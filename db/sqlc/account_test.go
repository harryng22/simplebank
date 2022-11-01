package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/harryng22/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
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

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	expectedAccount := CreateRandomAccount(t)
	actualAccount, err := testQueries.GetAccount(context.Background(), expectedAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, actualAccount)
	require.Equal(t, expectedAccount.ID, actualAccount.ID)
	require.Equal(t, expectedAccount.Owner, actualAccount.Owner)
	require.Equal(t, expectedAccount.Balance, actualAccount.Balance)
	require.Equal(t, expectedAccount.Currency, actualAccount.Currency)
	require.WithinDuration(t, expectedAccount.CreatedAt, actualAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	actualAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, actualAccount)
	require.Equal(t, account1.ID, actualAccount.ID)
	require.Equal(t, account1.Owner, actualAccount.Owner)
	require.Equal(t, arg.Balance, actualAccount.Balance)
	require.Equal(t, account1.Currency, actualAccount.Currency)
	require.WithinDuration(t, account1.CreatedAt, actualAccount.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestGetListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
