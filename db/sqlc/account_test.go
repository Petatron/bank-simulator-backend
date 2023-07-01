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

	Context("customer functions", func() {
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
