package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Petatron/bank-simulator-backend/db/mock"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/Petatron/bank-simulator-backend/db/util"
	"github.com/Petatron/bank-simulator-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestAccountAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for APIs")

}

func addAuthorizations(
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType, username string,
	duration time.Duration,
) {
	resultToken, err := tokenMaker.CreateToken(username, duration)
	if err != nil {
		panic(err)
	}
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, resultToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

var _ = Describe("API tests", func() {
	Context("getAccount API", func() {
		userName := util.GetRandomOwnerName()
		account := getRandomAccount(userName)

		testCases := []struct {
			name          string
			accountID     int64
			setupAuth     func(request *http.Request, tokenMaker token.Maker)
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "URI Binding Error",
				accountID: 0, // Assuming 0 is an invalid ID for testing
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					// No stubbing needed for this test
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},
			{
				name:      "Account Ownership Error",
				accountID: account.ID, // Assuming this is a valid account ID
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					// Set up authorization for a user different from the account owner
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, "otherUser", time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					// Stub the database to return an account owned by a different user
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(db.Account{Owner: "actualOwner", ID: account.ID}, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
				},
			},
			{
				name:      "Unauthorized User",
				accountID: account.ID,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, "unauthorized", time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(account, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
				},
			},

			{
				name:      "No Authorization",
				accountID: account.ID,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
				},
			},

			{
				name:      "Not Found",
				accountID: account.ID,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
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
				name:      "Account Not Found",
				accountID: account.ID, // Assuming this is a non-existent account ID
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
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
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
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
				name:      "Internal Server Error",
				accountID: account.ID, // Can use a valid or invalid ID here; it won't matter because the stub will force an error
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(account.ID)).
						Times(1).
						Return(db.Account{}, errors.New("internal server error"))
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},
			{
				name:      "Invalid ID",
				accountID: 0,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},
			{
				name:      "getAccount OK",
				accountID: account.ID,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
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
				server := newTestServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/accounts/%d", tc.accountID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				Expect(err).ShouldNot(HaveOccurred())

				tc.setupAuth(request, server.tokenMaker)

				// call the server
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.checkResponse(recorder)
			})
		}
	})

	Context("listAccounts API", func() {
		userName := util.GetRandomOwnerName()
		accounts := make([]db.Account, 5)
		for i := range accounts {
			accounts[i] = getRandomAccount(userName)
		}
		pageSize := int32(5)

		type Query struct {
			PageID   int32
			PageSize int32
		}

		testCases := []struct {
			name          string
			setupAuth     func(request *http.Request, tokenMaker token.Maker)
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
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListAccountsParams{
						Owner:  userName,
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
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
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
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListAccountsParams{
						Owner:  userName,
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
				server := newTestServer(store)
				recorder := httptest.NewRecorder()

				url := "/accounts"

				request, err := http.NewRequest(http.MethodGet, url, nil)
				Expect(err).ShouldNot(HaveOccurred())

				q := request.URL.Query()
				q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
				q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
				request.URL.RawQuery = q.Encode()

				tc.setupAuth(request, server.tokenMaker)

				// call the server
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.checkResponse(recorder)
			})
		}

	})

	Context("createAccount API", func() {
		userName := util.GetRandomOwnerName()
		account := getRandomAccount(userName)

		testCases := []struct {
			name          string
			body          gin.H
			accountID     int64
			setupAuth     func(request *http.Request, tokenMaker token.Maker)
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name: "OK",
				body: gin.H{
					"owner":    account.Owner,
					"currency": account.Currency,
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
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
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
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
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
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
				name: "Internal Error",
				body: gin.H{
					"owner":    account.Owner,
					"currency": account.Currency,
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
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

			{
				name: "Status Forbidden",
				body: gin.H{
					"owner":    account.Owner,
					"currency": account.Currency,
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0,
					}

					var pqError *pq.Error
					pqError = &pq.Error{
						Code: "23505",
					}

					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(db.Account{}, pqError)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusForbidden))
				},
			},
			{
				name: "Foreign Key Violation",
				body: gin.H{
					"owner":    "user1",
					"currency": "USD",
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, "user1", time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    "user1",
						Currency: "USD",
						Balance:  0,
					}
					pqError := &pq.Error{Code: "23503"} // assuming 23503 is the code for a foreign key violation
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(db.Account{}, pqError)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusForbidden))
				},
			},
			{
				name: "Unique Violation Error",
				body: gin.H{
					"owner":    userName, // assuming the username is already taken
					"currency": "USD",
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    userName,
						Currency: "USD",
						Balance:  0,
					}
					pqError := &pq.Error{
						Code: "23505", // PostgreSQL error code for unique_violation
					}
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(db.Account{}, pqError)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusForbidden))
				},
			},
			{
				name: "Internal Server Error",
				body: gin.H{
					"owner":    userName,
					"currency": "USD",
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    userName,
						Currency: "USD",
						Balance:  0,
					}
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(db.Account{}, errors.New("unexpected database error"))
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},
			{
				name: "Successful Account Creation",
				body: gin.H{
					"currency": "USD",
				},
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner:    userName, // Assuming userName is obtained from the token
						Currency: "USD",
						Balance:  0,
					}
					// Mock account returned from the database operation
					mockedAccount := db.Account{
						ID:       1, // Assuming a successful creation returns an account with ID 1
						Owner:    userName,
						Balance:  0,
						Currency: "USD",
					}
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(mockedAccount, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusOK))

					var account db.Account
					err := json.NewDecoder(recorder.Body).Decode(&account)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(account.ID).To(Equal(int64(1)))
					Expect(account.Owner).To(Equal(userName))
					Expect(account.Currency).To(Equal("USD"))
					Expect(account.Balance).To(Equal(int64(0)))
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
				server := newTestServer(store)
				recorder := httptest.NewRecorder()

				body, err := json.Marshal(tc.body)
				Expect(err).ShouldNot(HaveOccurred())

				url := "/accounts"
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
				Expect(err).ShouldNot(HaveOccurred())

				tc.setupAuth(request, server.tokenMaker)

				// call the server
				server.router.ServeHTTP(recorder, request)
				// check the response
				tc.checkResponse(recorder)
			})
		}
	})

	Context("deleteAccount API", func() {
		userName := util.GetRandomOwnerName()
		account := getRandomAccount(userName)

		testCases := []struct {
			name          string
			accountID     int64
			setupAuth     func(request *http.Request, tokenMaker token.Maker)
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "URI Binding Error",
				accountID: 0, // Assuming 0 is an invalid ID to trigger the binding error
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					// No stubbing needed for this test
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},
			{
				name:      "Internal Server Error",
				accountID: account.ID,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
						Return(errors.New("internal error"))
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},
			{
				name:      "OK",
				accountID: account.ID,
				setupAuth: func(request *http.Request, tokenMaker token.Maker) {
					addAuthorizations(request, tokenMaker, authorizationTypeBearer, userName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
						Return(nil) // Simulate successful deletion
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				},
			},
		}

		for _, tc := range testCases {
			It(fmt.Sprintf("Test case: %s", tc.name), func() {
				// Create a mock store
				controller := gomock.NewController(GinkgoT())
				defer controller.Finish()
				store := mockdb.NewMockStore(controller)

				tc.buildStubs(store) // Setup the expected database interactions

				// Start test server and send request
				server := newTestServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/accounts/%d", tc.accountID)
				request, err := http.NewRequest(http.MethodDelete, url, nil)
				Expect(err).NotTo(HaveOccurred())

				tc.setupAuth(request, server.tokenMaker)

				// Call the server
				server.router.ServeHTTP(recorder, request)

				// Check the response
				tc.checkResponse(recorder)
			})
		}
	})

})

func getRandomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.GetRandomInt(),
		Owner:    owner,
		Balance:  util.GetRandomMoneyAmount(),
		Currency: util.GetRandomCurrency(),
	}
}

func requireBodyMatchAccount(body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	Expect(err).ShouldNot(HaveOccurred())

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(gotAccount).Should(Equal(account))
}

func requireBodyMatchAccounts(body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	Expect(err).ShouldNot(HaveOccurred())

	var gotAccount []db.Account
	err = json.Unmarshal(data, &gotAccount)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(gotAccount).Should(Equal(accounts))
}
