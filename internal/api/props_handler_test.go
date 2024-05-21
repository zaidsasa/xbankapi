package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	storageMocks "github.com/zaidsasa/xbankapi/internal/storage/mocks"
)

func TestPropsHandler_health(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	connMock := storageMocks.NewMockDBConnection(t)

	propsHandler := NewPropsHandler(connMock)

	propsHandler.health(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	got, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, responseOK, string(got))
}

func TestPropsHandler_readiness(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		mock           func(*storageMocks.MockDBConnection)
		wantStatusCode int
		want           string
	}{
		{
			name: "failed when database is not available",
			mock: func(md *storageMocks.MockDBConnection) {
				md.EXPECT().Ping(mock.Anything).Return(errAnything).Once()
			},
			wantStatusCode: http.StatusInternalServerError,
			want: `not ready
`,
		},
		{
			name: "success when database is available",
			mock: func(md *storageMocks.MockDBConnection) {
				md.EXPECT().Ping(mock.Anything).Return(nil).Once()
			},
			wantStatusCode: http.StatusOK,
			want:           responseOK,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/readiness", nil)
			w := httptest.NewRecorder()

			connMock := storageMocks.NewMockDBConnection(t)

			propsHandler := NewPropsHandler(connMock)

			if tt.mock != nil {
				tt.mock(connMock)
			}

			propsHandler.readiness(w, r)

			res := w.Result()
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			defer res.Body.Close()

			got, err := io.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}
