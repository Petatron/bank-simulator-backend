package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/Petatron/bank-simulator-backend/db/mock"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/Petatron/bank-simulator-backend/db/util"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for APIs")

}

var _ = Describe("API tests", func() {
	Context("getAccount API", func() {
		account := getRandomAccount()

		testCases := []struct {
			name          string
			accountID     int64
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "getAccount OK",
				accountID: account.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(account, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					requireBodyMatchAccount(recorder.Body, account)
					Expect(recorder.Code).To(Equal(http.StatusOK))
				},
			},

			{
				name:      "Not Found",
				accountID: account.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(db.Account{}, sql.ErrNoRows)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				},
			},

			{
				name:      "Internal Error",
				accountID: account.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},

			{
				name:      "Invalid ID",
				accountID: 0,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
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

				url := fmt.Sprintf("/accounts/%d", tc.accountID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				Expect(err).ShouldNot(HaveOccurred())

				// call the server
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.checkResponse(recorder)
			})
		}
	})

	Context("listAccounts API", func() {
		accounts := make([]db.Account, 5)
		for i := range accounts {
			accounts[i] = getRandomAccount()
		}
		pageSize := int32(5)

		type Query struct {
			PageID   int32
			PageSize int32
		}

		testCases := []struct {
			name          string
			accountID     int64
			query         Query
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name: "OK",
				query: Query{
					PageID:   1,
					PageSize: 5,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListAccountsParams{
						Limit:  pageSize,
						Offset: 0,
					}
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(accounts, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					requireBodyMatchAccounts(recorder.Body, accounts)
					Expect(recorder.Code).To(Equal(http.StatusOK))
				},
			},

			{
				name: "Bad Request",
				query: Query{
					PageID:   1,
					PageSize: 1,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListAccountsParams{
						Limit:  1,
						Offset: 0,
					}
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Eq(arg)).
						Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},

			{
				name: "Internal Error",
				query: Query{
					PageID:   1,
					PageSize: 5,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListAccountsParams{
						Limit:  pageSize,
						Offset: 0,
					}
					store.EXPECT().
						ListAccounts(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return([]db.Account{}, sql.ErrConnDone)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
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

				url := "/accounts"

				request, err := http.NewRequest(http.MethodGet, url, nil)
				Expect(err).ShouldNot(HaveOccurred())

				q := request.URL.Query()
				q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
				q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
				request.URL.RawQuery = q.Encode()

				// call the server
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.checkResponse(recorder)
			})
		}

	})

	Context("createAccount API", func() {
		account := getRandomAccount()
		account.Owner = util.GetRandomOwnerName()
		account.Currency = util.GetRandomCurrency()

		testCases := []struct {
			name          string
			body          gin.H
			accountID     int64
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name: "OK",
				body: gin.H{
					"owner":    account.Owner,
					"currency": account.Currency,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0,
					}
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(account, nil)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					requireBodyMatchAccount(recorder.Body, account)
					Expect(recorder.Code).To(Equal(http.StatusOK))
				},
			},

			{
				name: "Bad Request",
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},

			{
				name: "Not Valid Currency",
				body: gin.H{
					"owner":    account.Owner,
					"currency": "invalid",
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0,
					}
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(0)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},

			{
				name: "Internal Error",
				body: gin.H{
					"owner":    account.Owner,
					"currency": account.Currency,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0,
					}
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
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

				url := "/accounts"
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

func getRandomAccount() db.Account {
	return db.Account{
		ID:       util.GetRandomInt(),
		Owner:    util.GetRandomOwnerName(),
		Balance:  util.GetRandomMoneyAmount(),
		Currency: util.GetRandomCurrency(),
	}
}

func requireBodyMatchAccount(body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	Expect(err).ShouldNot(HaveOccurred())

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(gotAccount).Should(Equal(account))
}

func requireBodyMatchAccounts(body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	Expect(err).ShouldNot(HaveOccurred())

	var gotAccount []db.Account
	err = json.Unmarshal(data, &gotAccount)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(gotAccount).Should(Equal(accounts))
}
