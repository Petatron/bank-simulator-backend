package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/Petatron/bank-simulator-backend/db/mock"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("API tests", func() {
	Context("createTransfer API", func() {

		fromAccount := getRandomAccount()
		toAccount := getRandomAccount()

		testCases := []struct {
			name          string
			body          gin.H
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(recorder *httptest.ResponseRecorder)
		}{
			{
				name: "OK",

				body: gin.H{
					"from_account_id": fromAccount.ID,
					"to_account_id":   toAccount.ID,
					"amount":          10,
					"currency":        "USD",
				},

				buildStubs: func(store *mockdb.MockStore) {
					fromAccount.Currency = "USD"
					toAccount.Currency = "USD"

					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
						Times(1).
						Return(fromAccount, nil)
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
						Times(1).
						Return(toAccount, nil)
					arg := db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        10,
					}
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(arg)).
						Times(1)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				},
			},

			{
				name: "Bad Request",

				buildStubs: func(store *mockdb.MockStore) {
					fromAccount.Currency = "USD"
					toAccount.Currency = "USD"
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Any()).
						Times(0)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				},
			},

			{
				name: "Internal Server Error",

				body: gin.H{
					"from_account_id": fromAccount.ID,
					"to_account_id":   toAccount.ID,
					"amount":          10,
					"currency":        "USD",
				},

				buildStubs: func(store *mockdb.MockStore) {
					fromAccount.Currency = "USD"
					toAccount.Currency = "USD"

					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
						Times(1).
						Return(fromAccount, nil)
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
						Times(1).
						Return(toAccount, nil)
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.TransferTxResult{}, sql.ErrConnDone)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				},
			},

			{
				name: "From Account Not Found",

				body: gin.H{
					"from_account_id": fromAccount.ID,
					"to_account_id":   toAccount.ID,
					"amount":          10,
					"currency":        "USD",
				},

				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
						Times(1).
						Return(db.Account{}, sql.ErrNoRows)
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
						Times(0)
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Any()).
						Times(0)
				},

				checkResponse: func(recorder *httptest.ResponseRecorder) {
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				},
			},

			{
				name: "Account Currency Mismatch",

				body: gin.H{
					"from_account_id": fromAccount.ID,
					"to_account_id":   toAccount.ID,
					"amount":          10,
					"currency":        "USD",
				},

				buildStubs: func(store *mockdb.MockStore) {
					fromAccount.Currency = "USD"
					toAccount.Currency = "EUR"

					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
						Times(1).
						Return(fromAccount, nil)
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
						Times(1).
						Return(toAccount, nil)
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Any()).
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

				body, err := json.Marshal(tc.body)
				Expect(err).ShouldNot(HaveOccurred())

				url := "/transfers"
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
