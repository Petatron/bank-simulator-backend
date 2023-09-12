package db

import (
	"context"

	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL Entry Operations", func() {

	Context("Entry Operations", func() {
		It("Test CreateEntry", func() {
			testOwnerName := createRandomUser()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName.Username,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, _ := testQueries.CreateAccount(context.Background(), arg)

			entryAmount := util.GetRandomMoneyAmount()

			entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
				AccountID: account.ID,
				Amount:    entryAmount,
			})
			Expect(err).To(BeNil())
			Expect(entry.Amount).To(Equal(entryAmount))
		})

		It("Test GetEntry", func() {
			testOwnerName := createRandomUser()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName.Username,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, _ := testQueries.CreateAccount(context.Background(), arg)

			entryAmount := util.GetRandomMoneyAmount()

			entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
				AccountID: account.ID,
				Amount:    entryAmount,
			})
			Expect(err).To(BeNil())
			Expect(entry.Amount).To(Equal(entryAmount))

			entry, err = testQueries.GetEntry(context.Background(), entry.ID)
			Expect(err).To(BeNil())
			Expect(entry.Amount).To(Equal(entryAmount))
		})

		It("Test ListEntries", func() {
			testOwnerName := createRandomUser()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName.Username,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, _ := testQueries.CreateAccount(context.Background(), arg)

			entryAmount := util.GetRandomMoneyAmount()

			_, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
				AccountID: account.ID,
				Amount:    entryAmount,
			})
			Expect(err).To(BeNil())

			_, err = testQueries.CreateEntry(context.Background(), CreateEntryParams{
				AccountID: account.ID,
				Amount:    entryAmount,
			})
			Expect(err).To(BeNil())

			entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
				AccountID: account.ID,
				Limit:     5,
				Offset:    0,
			})
			Expect(err).To(BeNil())
			Expect(len(entries)).To(Equal(2))
		})

	})

})
