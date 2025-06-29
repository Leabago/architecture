package redirect

import (
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"url-shortener/internal/lib/api"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string // Имя теста
		alias     string // Отправляемый alias
		url       string
		respError string // Какую ошибку мы должны получить?
		mockError error  // Ошибку, которую вернёт мок
	}{
		{
			name:  "ok",
			alias: "google",
			url:   "https://www.google.com",
		},
	}

	for _, tc := range cases {

		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()

			r.Get("/{alias}", New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// Проверяем, что статус ответа корректный
			require.Equal(t, tc.url, redirectedToURL)

		})

	}
}

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// Check the final URL after redirection.
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
