package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Password", func() {
	Context("Password operations", func() {
		It("Test HashPassword", func() {
			password := GetRandomStringWithLength(10)
			hashedPassword, err := HashPassword(password)
			Expect(err).To(BeNil())
			Expect(hashedPassword).NotTo(Equal(""))

			err = CheckPassword(password, hashedPassword)
			Expect(err).To(BeNil())

			wrongPassword := GetRandomStringWithLength(10)
			err = CheckPassword(wrongPassword, hashedPassword)
			Expect(err).NotTo(BeNil())

			longPassword := GetRandomStringWithLength(73)
			_, err = HashPassword(longPassword)
			Expect(err).NotTo(BeNil())
		})
	})
})
