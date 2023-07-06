package db

import (
	"context"
	"testing"

	"github.com/Petatron/bank-simulator-backend/db/util"
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
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
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
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
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
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, err := testQueries.CreateAccount(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(account.Balance).To(Equal(arg.Balance))
			Expect(account.Currency).To(Equal(arg.Currency))
			Expect(account.Owner).To(Equal(arg.Owner))
			Expect(account.ID).NotTo(BeZero())
			Expect(account.CreatedAt).NotTo(BeZero())

			testAmountForUpdate := util.GetRandomMoneyAmount()
			updateArg := UpdateAccountParams{
				ID:      account.ID,
				Balance: testAmountForUpdate,
			}

			updateAccount, err := testQueries.UpdateAccount(context.Background(), updateArg)
			Expect(err).To(BeNil())
			Expect(updateAccount.Balance).To(Equal(updateArg.Balance))
		})
	})

	Context("SQL operations", func() {
		It("Test ListAccounts", func() {
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg1 := ListAccountsParams{
				Owner:  testOwnerName,
				Limit:  2,
				Offset: 0,
			}

			arg2 := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
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
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			testAccount := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, err := testQueries.CreateAccount(context.Background(), testAccount)
			Expect(err).To(BeNil())
			Expect(account).NotTo(BeNil())

			err = testQueries.DeleteAccount(context.Background(), account.ID)
			Expect(err).To(BeNil())

			_, err = testQueries.GetAccount(context.Background(), account.ID)
			Expect(err).NotTo(BeNil())

		})
	})

	Context("SQL operations", func() {
		It("Test GetAccountForUpdate", func() {
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			arg := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
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

	Context("SQL operations", func() {
		It("Test AddAccountBalance", func() {
			testOwnerName := util.GetRandomOwnerName()
			testBalance := util.GetRandomMoneyAmount()
			testCurrency := util.GetRandomCurrency()
			testAccount := CreateAccountParams{
				Owner:    testOwnerName,
				Balance:  testBalance,
				Currency: testCurrency,
			}

			account, err := testQueries.CreateAccount(context.Background(), testAccount)
			Expect(err).To(BeNil())
			Expect(account.Balance).To(Equal(testAccount.Balance))
			Expect(account.Currency).To(Equal(testAccount.Currency))
			Expect(account.Owner).To(Equal(testAccount.Owner))
			Expect(account.ID).NotTo(BeZero())
			Expect(account.CreatedAt).NotTo(BeZero())

			testAmountForUpdate := util.GetRandomMoneyAmount()
			addAccountBalanceArg := AddAccountBalanceParams{
				ID:     account.ID,
				Amount: testAmountForUpdate,
			}

			updatedAccount, err := testQueries.AddAccountBalance(context.Background(), addAccountBalanceArg)
			Expect(err).To(BeNil())
			Expect(updatedAccount.Balance).To(Equal(testAccount.Balance + addAccountBalanceArg.Amount))
		})
	})
})
