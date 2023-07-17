package util

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for Random Util Operations")
}

var _ = Describe("Operation", func() {

	BeforeEach(func() {})

	AfterEach(func() {})

	Context("Util operations", func() {
		It("Test GetRandomInt", func() {
			testNumber := GetRandomInt()
			Expect(testNumber).NotTo(BeNil())
			Expect(testNumber).NotTo(Equal(0))
		})
	})

	Context("Util operations", func() {
		It("Test GetRandomIntWithRange", func() {
			testNumber := GetRandomIntWithRange(1, 99)
			Expect(testNumber).To(BeNumerically(">=", 1))
			Expect(testNumber).To(BeNumerically("<=", 99))
		})
	})

	Context("Util operations", func() {
		It("Test GetRandomStringWithLength", func() {
			testString := GetRandomStringWithLength(5)
			Expect(testString).NotTo(Equal(""))
			Expect(len(testString)).To(Equal(5))
		})
	})

	Context("Util operations", func() {
		It("Test GetRandomName", func() {
			testRandomOwnerName := GetRandomOwnerName()
			Expect(testRandomOwnerName).NotTo(BeNil())
			Expect(testRandomOwnerName).NotTo(Equal(""))
		})
	})

	Context("Util operations", func() {
		It("Test GetRandomCurrency", func() {
			testCurrency := GetRandomCurrency()
			Expect(testCurrency).NotTo(BeNil())
			Expect(testCurrency).NotTo(Equal(""))

			random := GetRandomInt()
			Expect(random).NotTo(BeNil())
		})
	})
})
