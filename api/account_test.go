package api

import (
	"bytes"
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
			controller := gomock.NewController(GinkgoT())
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", account.ID)
			request, err := http.NewRequest("GET", url, nil)
			Expect(err).ShouldNot(HaveOccurred())

			// call the server
			server.router.ServeHTTP(recorder, request)
			// check the response
			Expect(recorder.Code).To(Equal(http.StatusOK))
			requireBodyMatchAccount(recorder.Body, account)
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
