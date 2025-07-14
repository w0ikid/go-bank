package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/w0ikid/go-bank/db/mock"
	db "github.com/w0ikid/go-bank/db/sqlc"
)

func TestCreateTransferAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()

	tests := []struct {
		name          string
		amount        int64
		fromAccountID int64
		toAccountID   int64
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
		{
			name:          "OK",
			amount:        100,
			fromAccountID: account1.ID,
			toAccountID:   account2.ID,
			buildStubs: func(store *mock.MockStore) {
				account1.Currency = account2.Currency // Ensure currency match
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)
				store.EXPECT().
					CreateTransfer(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Transfer{
						ID:            1,
						FromAccountID: account1.ID,
						ToAccountID:   account2.ID,
						Amount:        100,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var transfer db.Transfer
				err := json.Unmarshal(recorder.Body.Bytes(), &transfer)
				require.NoError(t, err)
				require.Equal(t, account1.ID, transfer.FromAccountID)
				require.Equal(t, account2.ID, transfer.ToAccountID)
				require.Equal(t, int64(100), transfer.Amount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/transfers"
			requestBody := fmt.Sprintf(`{"from_account_id": %d, "to_account_id": %d, "amount": %d}`, tt.fromAccountID, tt.toAccountID, tt.amount)
			request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(requestBody))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")
			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}