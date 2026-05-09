package config

import (
	"net/http"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")

}
