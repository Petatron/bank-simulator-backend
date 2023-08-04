package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config Test", func() {
	It("Test LoadConfig with invalid path", func() {
		_, err2 := LoadConfig("")
		Expect(err2).To(Not(BeNil()))
	})

	It("Test LoadConfig", func() {
		config, err := LoadConfig("../..")
		Expect(err).To(BeNil())
		Expect(config.DBDriver).To(Equal("postgres"))
		Expect(config.DBSource).To(Equal("postgresql://root:password@localhost:5432/bank_simulator?sslmode=disable"))
		Expect(config.ServerAddress).To(Equal("0.0.0.0:8080"))
	})

})
