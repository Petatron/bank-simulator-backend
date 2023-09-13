package db

import (
	"context"
	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL Transfer Operations", func() {

	Context("Transfer Operations", func() {
		It("Test CreateTransfer", func() {
			testOwnerName := createRandomUser()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName.Username,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, _ := testQueries.CreateAccount(context.Background(), arg)

			transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
				FromAccountID: account.ID,
				ToAccountID:   account.ID,
				Amount:        100,
			})
			Expect(err).To(BeNil())
			Expect(transfer.Amount).To(Equal(int64(100)))
		})

		It("Test GetTransfer", func() {
			testOwnerName := createRandomUser()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName.Username,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, _ := testQueries.CreateAccount(context.Background(), arg)

			transferAmount := util.GetRandomMoneyAmount()

			transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
				FromAccountID: account.ID,
				ToAccountID:   account.ID,
				Amount:        transferAmount,
			})
			Expect(err).To(BeNil())
			Expect(transfer.Amount).To(Equal(transferAmount))
		})

		It("Test ListTransfers", func() {
			testOwnerName := createRandomUser()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName.Username,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, _ := testQueries.CreateAccount(context.Background(), arg)
			transferAmount := util.GetRandomMoneyAmount()

			_, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
				FromAccountID: account.ID,
				ToAccountID:   account.ID,
				Amount:        transferAmount,
			})
			Expect(err).To(BeNil())

			_, err = testQueries.CreateTransfer(context.Background(), CreateTransferParams{
				FromAccountID: account.ID,
				ToAccountID:   account.ID,
				Amount:        transferAmount,
			})
			Expect(err).To(BeNil())

			transfers, err := testQueries.ListTransfers(context.Background(), ListTransfersParams{
				FromAccountID: account.ID,
				ToAccountID:   account.ID,
				Limit:         5,
				Offset:        0,
			})
			Expect(err).To(BeNil())
			Expect(len(transfers)).To(Equal(2))

		})

	})

})
