package token

import (
	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for currency type operations")
}

var _ = Describe("PASETO maker tests", func() {
	It("Test create and verify token", func() {
		maker, err := NewPasetoMaker(util.GetRandomStringWithLength(32))
		Expect(err).To(BeNil())

		username := util.GetRandomOwnerName()
		duration := time.Minute

		issuedAt := time.Now()
		expiredAt := issuedAt.Add(duration)

		token, err := maker.CreateToken(username, duration)
		Expect(err).To(BeNil())
		Expect(token).NotTo(BeEmpty())

		payload, err := maker.VerifyToken(token)
		Expect(err).To(BeNil())
		Expect(payload).NotTo(BeNil())
		Expect(payload.Username).To(Equal(username))
		Expect(payload.IssuedAt).To(BeTemporally("~", issuedAt))
		Expect(payload.ExpiredAt).To(BeTemporally("~", expiredAt))
	})

	It("Test expired token", func() {
		maker, err := NewPasetoMaker(util.GetRandomStringWithLength(32))
		Expect(err).To(BeNil())

		token, err := maker.CreateToken(util.GetRandomOwnerName(), -time.Minute)
		Expect(err).To(BeNil())
		Expect(token).NotTo(BeEmpty())

		payload, err := maker.VerifyToken(token)
		Expect(err).To(Equal(ErrExpiredToken))
		Expect(payload).To(BeNil())
	})
})
