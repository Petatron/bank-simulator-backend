package db

import (
	"context"
	"fmt"

	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createRandomAccount() Account {
	testOwnerName := createRandomUser()
	testBalance := util.GetRandomMoneyAmount()
	testCurrency := util.GetRandomCurrency()
	arg := CreateAccountParams{
		Owner:    testOwnerName.Username,
		Balance:  testBalance,
		Currency: testCurrency,
	}

	account, _ := testQueries.CreateAccount(context.Background(), arg)
	return account
}

var _ = Describe("Operation", func() {
	Context("DB operations", func() {
		It("Test DB operations", func() {
			store := NewStore(testDB)
			account1 := createRandomAccount()
			account2 := createRandomAccount()
			fmt.Println(">>> Before transfer, amount for two accounts: ", account1.Balance, account2.Balance)

			n := 5
			amount := int64(10)

			errs := make(chan error)
			results := make(chan TransferTxResult)

			for i := 0; i < n; i++ {
				go func() {
					result, err := store.TransferTx(context.Background(), TransferTxParams{
						FromAccountID: account1.ID,
						ToAccountID:   account2.ID,
						Amount:        amount,
					})

					errs <- err
					results <- result
				}()
			}

			existed := make(map[int]bool)
			// Results check
			for i := 0; i < n; i++ {
				err := <-errs
				result := <-results
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())

				transfer := result.Transfer

				// Check transfer
				Expect(transfer).NotTo(BeNil())
				Expect(transfer.FromAccountID).To(Equal(account1.ID))
				Expect(transfer.ToAccountID).To(Equal(account2.ID))
				Expect(transfer.ID).NotTo(Equal(0))
				Expect(transfer.CreatedAt).NotTo(Equal(0))
				Expect(amount).To(Equal(transfer.Amount))

				_, err = store.GetTransfer(context.Background(), transfer.ID)
				Expect(err).To(BeNil())

				// Check Entries
				fromEntry := result.FromEntry
				Expect(fromEntry).NotTo(BeNil())
				Expect(fromEntry.AccountID).To(Equal(account1.ID))
				Expect(fromEntry.Amount).To(Equal(-amount))
				Expect(fromEntry.ID).NotTo(Equal(0))
				Expect(fromEntry.CreatedAt).NotTo(Equal(0))

				_, err = store.GetEntry(context.Background(), fromEntry.ID)
				Expect(err).To(BeNil())

				toEntry := result.ToEntry
				Expect(toEntry).NotTo(BeNil())
				Expect(toEntry.AccountID).To(Equal(account2.ID))
				Expect(toEntry.Amount).To(Equal(amount))
				Expect(toEntry.ID).NotTo(Equal(0))
				Expect(toEntry.CreatedAt).NotTo(Equal(0))

				_, err = store.GetEntry(context.Background(), toEntry.ID)
				Expect(err).To(BeNil())

				// Check accounts
				fromAccount := result.FromAccount
				Expect(fromAccount).NotTo(BeNil())
				Expect(fromAccount.ID).To(Equal(account1.ID))

				toAccount := result.ToAccount
				Expect(toAccount).NotTo(BeNil())
				Expect(toAccount.ID).To(Equal(account2.ID))

				// Check account balance
				fmt.Println(">>> After transfer, amount for two accounts: ", fromAccount.Balance, toAccount.Balance)
				diff1 := account1.Balance - fromAccount.Balance
				diff2 := toAccount.Balance - account2.Balance
				Expect(diff1).To(Equal(diff2))
				Expect(diff1).To(BeNumerically(">", 0))
				Expect(diff1 % amount).To(Equal(int64(0)))

				k := int(diff1 / amount)
				Expect(k >= 1 && k <= n).To(BeTrue())
				Expect(existed[k]).To(BeFalse())
				existed[k] = true
			}

			// Check the final account balance
			updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
			Expect(err).To(BeNil())

			updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
			Expect(err).To(BeNil())

			fmt.Println(">>> Final amount for two accounts: ", updateAccount1.Balance, updateAccount2.Balance)
			Expect(updateAccount1.Balance).To(Equal(account1.Balance - int64(n)*amount))
			Expect(updateAccount2.Balance).To(Equal(account2.Balance + int64(n)*amount))
		})

		It("Test DB operations deadlock", func() {
			store := NewStore(testDB)
			account1 := createRandomAccount()
			account2 := createRandomAccount()
			fmt.Println(">>> Before transfer, amount for two accounts: ", account1.Balance, account2.Balance)

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
					_, err := store.TransferTx(context.Background(), TransferTxParams{
						FromAccountID: fromAccountID,
						ToAccountID:   toAccountID,
						Amount:        amount,
					})

					errs <- err
				}()
			}

			// Results check
			for i := 0; i < n; i++ {
				err := <-errs
				Expect(err).To(BeNil())
			}

			// Check the final account balance
			updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
			Expect(err).To(BeNil())

			updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
			Expect(err).To(BeNil())

			fmt.Println(">>> Final amount for two accounts: ", updateAccount1.Balance, updateAccount2.Balance)
			Expect(updateAccount1.Balance).To(Equal(account1.Balance))
			Expect(updateAccount2.Balance).To(Equal(account2.Balance))
		})
	})
})
