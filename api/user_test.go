package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "github.com/Petatron/bank-simulator-backend/db/mock"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/Petatron/bank-simulator-backend/db/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"
)

var _ = Describe("API tests", func() {
	Context("createUser API", func() {
		password, user := randomUserWithPassword()

		testCases := []struct {
			name          string
			body          gin.H
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name: "OK",
				body: gin.H{
					"username":  user.Username,
					"password":  password,
					"full_name": user.FullName,
					"email":     user.Email,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateUsersParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}
					store.EXPECT().
						CreateUsers(gomock.Any(), EqCreateUserParams(arg, password)).
						Times(1).
						Return(user, nil)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					requireBodyMatchUser(recorder.Body, user)
					Expect(recorder.Code).To(Equal(http.StatusOK))
				},
			},

			{
				name: "Bad Request",
				body: gin.H{
					"username":  user.Username,
					"password":  password,
					"full_name": user.FullName,
					"email":     "",
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUsers(gomock.Any(), gomock.Any()).
						Times(0)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},

			{
				name: "Internal Server Error",
				body: gin.H{
					"username":  user.Username,
					"password":  password,
					"full_name": user.FullName,
					"email":     user.Email,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateUsersParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}
					store.EXPECT().
						CreateUsers(gomock.Any(), EqCreateUserParams(arg, password)).
						Times(1).
						Return(db.User{}, fmt.Errorf("some error"))
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},

			{
				name: "Internal Server Error long password",
				body: gin.H{
					"username":  user.Username,
					"password":  util.GetRandomStringWithLength(73),
					"full_name": user.FullName,
					"email":     user.Email,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateUsersParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}
					store.EXPECT().
						CreateUsers(gomock.Any(), EqCreateUserParams(arg, password)).
						Times(0)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},

			{
				name: "Username Already Exists",
				body: gin.H{
					"username":  user.Username,
					"password":  password,
					"full_name": user.FullName,
					"email":     user.Email,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateUsersParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					}

					var pqError *pq.Error
					pqError = &pq.Error{
						Code: "23505",
					}

					store.EXPECT().
						CreateUsers(gomock.Any(), EqCreateUserParams(arg, password)).
						Times(1).
						Return(db.User{}, pqError)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusForbidden))
				},
			},
		}
		for i := range testCases {
			tc := testCases[i]

			It(fmt.Sprintf("Test case #%d: %s", i, tc.name), func() {
				// create mock store
				controller := gomock.NewController(GinkgoT())
				defer controller.Finish()

				store := mockdb.NewMockStore(controller)
				tc.buildStubs(store)

				// start test server and send request
				server := NewServer(store)
				recorder := httptest.NewRecorder()

				body, err := json.Marshal(tc.body)
				Expect(err).ShouldNot(HaveOccurred())

				url := "/users"
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
				Expect(err).ShouldNot(HaveOccurred())

				// call the server
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.checkResponse(recorder)
			})
		}
	})
})

func randomUserWithPassword() (password string, user db.User) {
	password, err := util.HashPassword(util.GetRandomStringWithLength(10))
	Expect(err).ShouldNot(HaveOccurred())
	return password,
		db.User{
			Username:          util.GetRandomOwnerName(),
			FullName:          util.GetRandomOwnerName(),
			Email:             util.GetRandomEmail(),
			PasswordChangedAt: time.Now(),
			CreatedAt:         time.Now(),
		}
}

func requireBodyMatchUser(body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	Expect(err).ShouldNot(HaveOccurred())

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	Expect(err).ShouldNot(HaveOccurred())

	Expect(gotUser.Username).Should(Equal(user.Username))
	Expect(gotUser.FullName).Should(Equal(user.FullName))
	Expect(gotUser.Email).Should(Equal(user.Email))
	Expect(gotUser.HashedPassword).Should(Equal(""))
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUsersParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUsersParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// EqCreateUserParams is a custom matcher that checks if two userParams are equal
func EqCreateUserParams(arg db.CreateUsersParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		arg,
		password,
	}
}
