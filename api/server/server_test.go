package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(SetupServer())

	defer ts.Close()

	response, err := http.Get(fmt.Sprintf("%s/health", ts.URL))

	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("expected status 200, got %v", response.StatusCode)
	}
}
