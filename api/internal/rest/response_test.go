package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/gin-gonic/gin"
)

func TestSetResponse(t *testing.T) {
	tests := []struct {
		err      error
		success  bool
		response ErrorResponse
	}{
		{
			err:     nil,
			success: true,
		},
		{
			success:  false,
			err:      errs.RestBadRequest("bad request"),
			response: ErrorResponse{Message: "bad request", Code: "bad_request"},
		},
		{
			success:  false,
			err:      errors.New("random err"),
			response: ErrorResponse{Message: "random err", Code: "internal_error"},
		},
		{
			success:  false,
			err:      errs.StoreNotFound("not found"),
			response: ErrorResponse{Message: "not found", Code: "not_found"},
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		var err error

		r.GET("/test", func(c *gin.Context) {
			SetResponse(c, struct{}{}, tc.err)
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		c.Request, err = http.NewRequest(http.MethodGet, "/test", nil)

		r.ServeHTTP(res, c.Request)

		if tc.success {
			if res.Result().StatusCode != 200 {
				t.Errorf("expected 200 status, got %v", res.Result().StatusCode)
			}
		} else {
			var response ErrorResponse

			err = json.NewDecoder(res.Body).Decode(&response)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if response.Code != tc.response.Code {
				t.Errorf("expected code '%v', but got '%v'", tc.response.Code, response.Code)
			}

			if response.Message != tc.response.Message {
				t.Errorf("expected message '%v', but got '%v'", tc.response.Message, response.Message)
			}
		}
	}
}
