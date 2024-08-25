package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func TestGetMediaAPI(t *testing.T) {
	media := randomMedia()

	testCases := []struct {
		name string
		mediaID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
		{
			name: "OK",
			mediaID: media.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				GetMedia(gomock.Any(),
				gomock.Eq(media.ID)).
				Times(1).
				Return(media, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMedia(t, recorder.Body, media)
			},
		},
		{
			name: "NotFound",
			mediaID: media.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().
				GetMedia(gomock.Any(), gomock.Eq(media.ID)).
				Times(1).
				Return(db.Medium{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/medium/%d", tc.mediaID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomMedia() db.Medium {
	return db.Medium{
		ID: util.RandomInt(1, 1000),
		Url: util.RandomString(10),
		AwsID: util.RandomString(15),
		MediaRef: util.RandomString(10),
	}
}

func requireBodyMatchMedia(t *testing.T, body *bytes.Buffer, media db.Medium) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMedium db.Medium
	err = json.Unmarshal(data, &gotMedium)
	require.NoError(t, err)
	require.Equal(t, media, gotMedium)
}