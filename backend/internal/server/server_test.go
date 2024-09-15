package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattjmcnaughton/go-youtube-feed/internal/youtube"
)

func TestGetStatus(t *testing.T) {
	// TODO: Come up w/ a better option for the mock youtube client - via some form of `interface`, etc.
	fakeYoutubeClient := youtube.NewYoutubeClient("fake-api-key")

	router := GetRouter(fakeYoutubeClient)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "OK")
}
