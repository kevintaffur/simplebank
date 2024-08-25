package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/kvgtl/simplebank/db/mock"
	db "github.com/kvgtl/simplebank/db/sqlc"
	"github.com/kvgtl/simplebank/utils"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	args     db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	args, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, args.HashedPassword)
	if err != nil {
		return false
	}

	e.args.HashedPassword = args.HashedPassword

	return reflect.DeepEqual(e.args, args)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches args %v and password %v", e.args, e.password)
}

func eqCreateUserParams(args db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{args, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
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
				args := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: user.HashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(args, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// this is the mocked version of the function that interact
				// with db, and the "Return()" part is what that function
				// returns, not the api response.
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(args, password)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			// this is what checks the api endpoint response.
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// this is the mocked version of the function that interact
				// with db, and the "Return()" part is what that function
				// returns, not the api response.
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(args, password)).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			// this is what checks the api endpoint response.
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "#3@ABCs@@@##",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// this is the mocked version of the function that interact
				// with db, and the "Return()" part is what that function
				// returns, not the api response.
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(args, password)).
					Times(0)
			},
			// this is what checks the api endpoint response.
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "email123",
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// this is the mocked version of the function that interact
				// with db, and the "Return()" part is what that function
				// returns, not the api response.
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(args, password)).
					Times(0)
			},
			// this is what checks the api endpoint response.
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				// this is the mocked version of the function that interact
				// with db, and the "Return()" part is what that function
				// returns, not the api response.
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(args, password)).
					Times(0)
			},
			// this is what checks the api endpoint response.
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		// run each test case as a subtest of this test.
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			testCase.buildStubs(store)

			// start test server and send requests.
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// marshall body data to json.
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))

			require.NoError(t, err)

			// run the server and record the response on recorder, the request was defined before.
			server.router.ServeHTTP(recorder, request)

			// check response.
			testCase.checkResponse(t, recorder)
		})

	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomString(10),
		Email:          utils.RandomEmailAddress(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	// Body contains the user.
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
