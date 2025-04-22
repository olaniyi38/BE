package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/olaniyi38/BE/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	randomName := user.Username
	randomBalance := util.RandomMoney()
	randomCurrency := util.RandomCurrency()

	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Name:     randomName,
		Balance:  randomBalance,
		Currency: randomCurrency,
	})
	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, randomName, account.Name)
	require.Equal(t, randomBalance, account.Balance)
	require.Equal(t, randomCurrency, account.Currency)
	require.NotZero(t, account.UpdatedAt)
	require.NotZero(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	gotAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)

	require.Equal(t, gotAccount, account)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:     account1.ID,
		Amount: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.UpdatedAt, account2.UpdatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
		Name:   lastAccount.Name,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Name, account.Name)
	}
}
