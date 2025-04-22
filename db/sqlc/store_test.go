package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	//we want to simulate multiple transfers to an account to test transactions properly

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	amount := int64(10)
	n := 5

	fmt.Println(">> before:", fromAccount.Balance, toAccount.Balance)

	results := make(chan TransferTxResult)
	errs := make(chan error)
	idempotentMap := map[int64]struct{}{}

	defer close(results)
	defer close(errs)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTX(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			results <- result
			errs <- err

		}()
	}

	for i := 1; i <= n; i++ {
		t.Logf("tx-%v >>>>>>>>>", i)
		t.Log(idempotentMap)

		//for each result
		result := <-results
		require.NotEmpty(t, result)

		err := <-errs
		require.NoError(t, err)
		t.Logf("amount %v \n", int64(i)*amount)
		t.Logf("from account balance og >>> %v, after >>> %v \n", fromAccount.Balance, result.FromAccount.Balance)
		t.Logf("to account balance og >>> %v, after >>> %v \n", toAccount.Balance, result.ToAccount.Balance)
		//check transfer
		require.NotZero(t, result.Transfer.ID)
		require.Equal(t, fromAccount.ID, result.Transfer.FromAccountID)
		require.Equal(t, toAccount.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.CreatedAt)

		//check entries
		require.NotZero(t, result.FromEntry.ID)
		require.Equal(t, result.FromEntry.AccountID, result.Transfer.FromAccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)

		require.NotZero(t, result.ToEntry.ID)
		require.Equal(t, result.ToEntry.AccountID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.ToEntry.Amount)

		//check accounts
		require.Equal(t, result.FromAccount.ID, fromAccount.ID)
		require.Equal(t, result.ToAccount.ID, toAccount.ID)

		//check that the correct amount entered and left each accounts
		fromAccountDiff := fromAccount.Balance - result.FromAccount.Balance
		toAccountDiff := result.ToAccount.Balance - toAccount.Balance
		require.Equal(t, fromAccountDiff, toAccountDiff)

		//check that each balance update is unique
		//ogBalance - current balance / amount === i
		//this should be unique for each iteration, 1st: 10 / 10 = 1, 2nd 20/10 =2, 3rd, 30/10 = 3
		idempotentKey := fromAccountDiff / amount
		t.Logf("idempoKey: %v \n", idempotentKey)
		require.EqualValues(t, idempotentKey, i)

		//check that the original balance each iteration has decreased or incremented by i * amount
		//i.e og bal = 40, 1st: 30, 2nd: 20, 3rd: 10
		//current balance === og balance - (i * amount) || currentBalance === og balance + (i * amount)
		require.Equal(t, result.FromAccount.Balance, fromAccount.Balance-(int64(i)*amount))
		require.Equal(t, result.ToAccount.Balance, toAccount.Balance+(int64(i)*amount))

		//check that this update is not running twice
		require.NotContains(t, idempotentMap, idempotentKey)
		idempotentMap[idempotentKey] = struct{}{}

	}

	// check the final updated balance
	updatedFromAccount, err := testStore.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testStore.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, fromAccount.Balance-(int64(n)*amount), updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+(int64(n)*amount), updatedToAccount.Balance)

}



func TestTransferTxDeadlock(t *testing.T) {
	//simulates multiple transactions between 2 accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testStore.TransferTX(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	//values are the same because they balance each other out
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}