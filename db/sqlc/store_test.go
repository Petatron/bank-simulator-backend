package db

import (
	"context"
	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/gomega"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	testOwnerName := util.GetRandomOwnerName()
	testBalance := util.GetRandomMoneyAmount()
	testCurrency := util.GetRandomCurrency()
	arg := CreateAccountParams{
		Owner:    testOwnerName,
		Balance:  testBalance,
		Currency: testCurrency,
	}

	account, _ := testQueries.CreateAccount(context.Background(), arg)
	return account
}

func TestDBTrans(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := util.GetRandomMoneyAmount()

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		result, err := store.TransferTx(context.Background(), TransferTxParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        amount,
		})
		if err != nil {
			t.Errorf("Error while creating entry: %v", err)
		}

		errs <- err
		results <- result
	}

	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		Expect(err).To(BeNil())
		Expect(result.FromAccount).To(Equal(amount))
	}
}
