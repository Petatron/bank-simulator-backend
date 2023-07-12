package db

import (
	"context"
	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
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
	RegisterFailHandler(Fail)
	defer GinkgoRecover()

	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := util.GetRandomMoneyAmount()

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
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
		}()
	}

	// Results check
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		Expect(err).To(BeNil())
		Expect(result).NotTo(BeNil())
		Expect(result.FromAccount).To(Equal(amount))
		Expect(result.ToAccount).To(Equal(amount))

		// Check transfer
		transfer := result.Transfer
		Expect(transfer).NotTo(BeNil())
		Expect(transfer.Amount).To(Equal(amount))
		Expect(transfer.FromAccountID).To(Equal(account1.ID))
		Expect(transfer.ToAccountID).To(Equal(account2.ID))
		Expect(transfer.ID).NotTo(Equal(0))
		Expect(transfer.CreatedAt).NotTo(Equal(0))
		Expect(amount).To(Equal(transfer.Amount))

		_, err = store.GetTransfer(context.Background(), account1.ID)
		Expect(err).To(BeNil())

		// Check Entries
		fromEntry := result.FromEntry
		Expect(fromEntry).NotTo(BeNil())
		Expect(fromEntry.AccountID).To(Equal(account1.ID))
		Expect(fromEntry.Amount).To(Equal(-amount))
		Expect(fromEntry.ID).NotTo(Equal(0))
		Expect(fromEntry.CreatedAt).NotTo(Equal(0))

		toEntry := result.ToEntry
		Expect(toEntry).NotTo(BeNil())
		Expect(toEntry.AccountID).To(Equal(account2.ID))
		Expect(toEntry.Amount).To(Equal(amount))
		Expect(toEntry.ID).NotTo(Equal(0))
		Expect(toEntry.CreatedAt).NotTo(Equal(0))

		_, err = store.GetEntry(context.Background(), account1.ID)
		Expect(err).To(BeNil())

	}
}

func TestDBT(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()

	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	//account2 := createRandomAccount(t)

	var amount = int64(-1000)
	_, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Amount:        amount,
	})

	if err != nil {
		t.Errorf("Error while creating entry: %v", err)
	}
}
