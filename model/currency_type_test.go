package model_test

import (
	"testing"

	"github.com/Petatron/bank-simulator-backend/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for currency type operations")
}

var _ = Describe("CurrencyType", func() {
	It("Test valid currency type", func() {
		Expect(model.USD.IsValid()).To(BeTrue())
		Expect(model.EUR.IsValid()).To(BeTrue())
		Expect(model.CAD.IsValid()).To(BeTrue())
	})

	It("Test invalid currency type", func() {
		Expect(model.CurrencyType("INVALID").IsValid()).To(BeFalse())
	})
})
