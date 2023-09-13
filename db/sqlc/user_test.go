package db

import (
	"context"
	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operation", func() {

	Context("SQL operations", func() {
		It("Test CreateUser", func() {
			userName := util.GetRandomOwnerName()
			hashedPassword := util.GetRandomStringWithLength(10)
			fullName := util.GetRandomOwnerName()
			email := util.GetRandomEmail()
			arg := CreateUsersParams{
				Username:       userName,
				HashedPassword: hashedPassword,
				FullName:       fullName,
				Email:          email,
			}

			user, err := testQueries.CreateUsers(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(user.Username).To(Equal(arg.Username))
		})

		It("Test GetUser", func() {

			userName := util.GetRandomOwnerName()
			hashedPassword := util.GetRandomStringWithLength(10)
			fullName := util.GetRandomOwnerName()
			email := util.GetRandomEmail()
			arg := CreateUsersParams{
				Username:       userName,
				HashedPassword: hashedPassword,
				FullName:       fullName,
				Email:          email,
			}

			user, err := testQueries.CreateUsers(context.Background(), arg)
			Expect(err).To(BeNil())
			Expect(user.Username).To(Equal(arg.Username))

			user, err = testQueries.GetUser(context.Background(), arg.Username)
			Expect(err).To(BeNil())
			Expect(user.Username).To(Equal(arg.Username))
		})
	})

})
