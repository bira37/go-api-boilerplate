package rest

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	response, err := http.Get(fmt.Sprintf("%s/health", Server.URL))

	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("expected status 200, got %v", response.StatusCode)
	}
}
