package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/Petatron/bank-simulator-backend/db/mock"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/Petatron/bank-simulator-backend/db/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit Test for APIs")

}

var _ = Describe("API tests", func() {
	Context("getAccount API", func() {
		It("Test getAccount API", func() {
			account := getRandomAccount()

			testCases := []struct {
				name          string
				accountID     int64
				buildStubs    func(store *mockdb.MockStore)
				checkResponse func(recorder *httptest.ResponseRecorder)
			}{
				{
					name:      "OK",
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
			}
		})
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
