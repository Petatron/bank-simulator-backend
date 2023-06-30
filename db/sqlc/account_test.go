package db

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
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

	})
})
