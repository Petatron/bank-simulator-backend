package db

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for DB Operations")
}

var _ = Describe("Operation", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	Context("SQL operations", func() {
		It("Test CreateAccountParams", func() {
			arg := CreateAccountParams{
				Owner:    "1",
				Balance:  100,
				Currency: "USD",
			}

			account, err := testQueries.CreateAccount(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(account.Balance).To(Equal(arg.Balance))
			Expect(account.Currency).To(Equal(arg.Currency))
			Expect(account.Owner).To(Equal(arg.Owner))
			Expect(account.ID).NotTo(BeZero())
			Expect(account.CreatedAt).NotTo(BeZero())
		})
	})

	Context("SQL operations", func() {
		It("Test GetAccount", func() {
			arg := CreateAccountParams{
				Owner:    "1",
				Balance:  100,
				Currency: "USD",
			}

			account, err := testQueries.CreateAccount(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(account.Balance).To(Equal(arg.Balance))
			Expect(account.Currency).To(Equal(arg.Currency))
			Expect(account.Owner).To(Equal(arg.Owner))
			Expect(account.ID).NotTo(BeZero())
			Expect(account.CreatedAt).NotTo(BeZero())

			getAccount, err := testQueries.GetAccount(context.Background(), account.ID)
			Expect(err).To(BeNil())
			Expect(getAccount.Balance).To(Equal(arg.Balance))
			Expect(getAccount.Currency).To(Equal(arg.Currency))
			Expect(getAccount.Owner).To(Equal(arg.Owner))
			Expect(getAccount.ID).NotTo(BeZero())
			Expect(getAccount.CreatedAt).NotTo(BeZero())
		})
	})

	Context("SQL operations", func() {
		It("Test UpdateAccount", func() {
			arg := CreateAccountParams{
				Owner:    "KGB",
				Balance:  100,
				Currency: "USD",
			}

			account, err := testQueries.CreateAccount(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(account.Balance).To(Equal(arg.Balance))
			Expect(account.Currency).To(Equal(arg.Currency))
			Expect(account.Owner).To(Equal(arg.Owner))
			Expect(account.ID).NotTo(BeZero())
			Expect(account.CreatedAt).NotTo(BeZero())

			updateArg := UpdateAccountParams{
				ID:      account.ID,
				Balance: 200,
			}

			updateAccount, err := testQueries.UpdateAccount(context.Background(), updateArg)
			Expect(err).To(BeNil())
			Expect(updateAccount.Balance).To(Equal(updateArg.Balance))
		})
	})

	Context("SQL operations", func() {
		It("Test ListAccounts", func() {
			arg1 := ListAccountsParams{
				Owner:  "CIA",
				Limit:  2,
				Offset: 0,
			}

			arg2 := CreateAccountParams{
				Owner:    "CIA",
				Balance:  100,
				Currency: "USD",
			}

			testAccount, err := testQueries.CreateAccount(context.Background(), arg2)
			Expect(err).To(BeNil())
			Expect(testAccount.Owner).To(Equal(arg2.Owner))

			accounts, err := testQueries.ListAccounts(context.Background(), arg1)
			Expect(err).To(BeNil())
			Expect(accounts).NotTo(BeNil())
			Expect(len(accounts)).To(BeNumerically(">=", 1))
		})

		It("Test DeleteAccountParams", func() {
			// Create a new account
			testAccount := CreateAccountParams{
				Owner:    "1",
				Balance:  100,
				Currency: "USD",
			}

			account, err := testQueries.CreateAccount(context.Background(), testAccount)
			Expect(err).To(BeNil())
			Expect(account).NotTo(BeNil())

			// Delete the account
			err = testQueries.DeleteAccount(context.Background(), account.ID)
			Expect(err).To(BeNil())

			// Getting deleted account should get error
			_, err = testQueries.GetAccount(context.Background(), account.ID)
			Expect(err).NotTo(BeNil())

		})
	})

	Context("SQL operations", func() {
		It("Test GetAccountForUpdate", func() {
			arg := CreateAccountParams{
				Owner:    "1",
				Balance:  100,
				Currency: "USD",
			}

			account, err := testQueries.CreateAccount(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(account.Balance).To(Equal(arg.Balance))
			Expect(account.Currency).To(Equal(arg.Currency))
			Expect(account.Owner).To(Equal(arg.Owner))
			Expect(account.ID).NotTo(BeZero())
			Expect(account.CreatedAt).NotTo(BeZero())

			getAccount, err := testQueries.GetAccountForUpdate(context.Background(), account.ID)
			Expect(err).To(BeNil())
			Expect(getAccount.Balance).To(Equal(arg.Balance))
			Expect(getAccount.Currency).To(Equal(arg.Currency))
			Expect(getAccount.Owner).To(Equal(arg.Owner))
			Expect(getAccount.ID).NotTo(BeZero())
			Expect(getAccount.CreatedAt).NotTo(BeZero())
		})
	})

})
