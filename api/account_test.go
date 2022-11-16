package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/harryng22/simplebank/db/mock"
	db "github.com/harryng22/simplebank/db/sqlc"
	"github.com/harryng22/simplebank/util"
	"github.com/stretchr/testify/require"
)

const (
	BASE_URL = "/accounts"
)

type AccountTestCase struct {
	name          string
	account       db.Account
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func GenerateGetAccountTestCases() []AccountTestCase {
	account := randomAccount()
	return []AccountTestCase{
		{
			name:    "OK",
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:    "NotFound",
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "InvalidID",
			account: db.Account{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
}

func GenerateCreateAccountTestCases() []AccountTestCase {
	account := randomAccount()
	args := db.CreateAccountParams{
		Owner:    account.Owner,
		Balance:  0,
		Currency: account.Currency,
	}
	return []AccountTestCase{
		{
			name:    "Ok",
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), args).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:    "BadRequest",
			account: db.Account{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), args).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "InternalServerError",
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), args).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
}

func TestGetAccountAPI(t *testing.T) {
	testCases := GenerateGetAccountTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/%d", BASE_URL, tc.account.ID)
			doTestAccountAPI(t, http.MethodGet, url, nil, tc)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	testCases := GenerateCreateAccountTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doTestAccountAPI(t, http.MethodPost, BASE_URL, tc.account, tc)
		})
	}
}

func doTestAccountAPI(t *testing.T, method string, url string, account any, tc AccountTestCase) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mockdb.NewMockStore(controller)

	// build stubs
	tc.buildStubs(store)

	// start test server and send request
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	body, err := json.Marshal(account)
	require.NoError(t, err)

	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	// check response
	tc.checkResponse(t, recorder)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
