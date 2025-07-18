package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/w0ikid/go-bank/db/mock"
	db "github.com/w0ikid/go-bank/db/sqlc"
	"github.com/w0ikid/go-bank/util"
	"reflect"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// сравниваем bcrypt
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	// сравниваем остальные поля
	e.arg.HashedPassword = arg.HashedPassword // обновляем на тот, что пришёл
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %s", e.arg, e.password)
}

// вспомогательная обёртка
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
    user, password := randomUser(t)

    tests := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mock.MockStore)
        checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
    }{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
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

            url := "/users"
            bodyBytes, err := json.Marshal(tt.body)
            require.NoError(t, err)

            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)

            tt.checkResponse(t, recorder)
        })
    }
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(10)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
}
