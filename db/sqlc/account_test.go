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

	//Context("SQL operations", func() {
	//	It("Test DeleteAccount", func() {
	//
	//		arg := CreateAccountParams{
	//			Owner:    "Steve",
	//			Balance:  100,
	//			Currency: "USD",
	//		}
	//
	//		account, err := testQueries.CreateAccount(context.Background(), arg)
	//		Expect(err).To(BeNil())
	//		Expect(account.Balance).To(Equal(arg.Balance))
	//		Expect(account.Currency).To(Equal(arg.Currency))
	//		Expect(account.Owner).To(Equal(arg.Owner))
	//		Expect(account.ID).NotTo(BeZero())
	//		Expect(account.CreatedAt).NotTo(BeZero())
	//
	//		err = testQueries.DeleteAccount(context.Background(), account.ID)
	//		Expect(err).To(BeNil())
	//
	//		getAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	//		Expect(err).ToNot(BeNil())
	//		Expect(getAccount).To(BeNil())
	//	})
	//})

	Context("SQL operations", func() {
		FIt("Test ListAccounts", func() {
			arg := CreateAccountParams{
				Owner:    "Steve",
				Balance:  100,
				Currency: "USD",
			}

			_, err := testQueries.CreateAccount(context.Background(), arg)
			Expect(err).To(BeNil())

			accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{})
			Expect(err).To(BeNil())
			Expect(accounts).To(HaveLen(1))
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

			deleteAccountId := account.ID

			// Delete the account
			err = testQueries.DeleteAccount(context.Background(), deleteAccountId)
			Expect(err).To(BeNil())

			// Getting deleted account should get error
			_, err = testQueries.GetAccount(context.Background(), deleteAccountId)
			Expect(err).NotTo(BeNil())

		})
	})

})
