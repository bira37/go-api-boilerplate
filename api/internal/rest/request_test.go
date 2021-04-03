package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/gin-gonic/gin"
)

type Object struct {
	Id    int     `binding:"required"`
	Text  string  `binding:"required"`
	Extra float64 `binding:"required"`
}

func TestParseBody(t *testing.T) {
	tests := []struct {
		request    interface{}
		success    bool
		err        *errs.RestError
		errorRegex string
	}{
		{
			request: Object{Id: 1, Text: "test", Extra: 2.0},
			success: true,
			err:     nil,
		},
		{
			request: struct {
				Id   string
				Text string
			}{Id: "1", Text: "test"},
			success:    false,
			err:        errs.RestBadRequest(""),
			errorRegex: "should have type",
		},
		{
			request: struct {
				Id   int
				Text string
			}{Id: 1, Text: "test"},
			success:    false,
			err:        errs.RestBadRequest(""),
			errorRegex: "'Extra' is required",
		},
		{
			request: struct {
				Id int
			}{Id: 1},
			success:    false,
			err:        errs.RestBadRequest(""),
			errorRegex: "Several errors occurred",
		},
	}

	for _, tc := range tests {
		res := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(res)

		var testerr error
		var err error

		r.POST("/test", func(c *gin.Context) {
			var obj Object
			testerr = ParseBody(c, &obj)
			c.JSON(200, struct{}{})
		})

		body, err := json.Marshal(tc.request)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		c.Request, err = http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))

		r.ServeHTTP(res, c.Request)

		if tc.success {
			if testerr != nil {
				t.Errorf("unexpected error: %v", testerr)
			}
		} else {
			if testerr == nil {
				t.Errorf("expected binding error: %v", testerr)
			}

			matched, err := regexp.MatchString(tc.errorRegex, testerr.Error())

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !matched {
				t.Errorf("expected error '%v' to match regex '%v'", testerr.Error(), tc.errorRegex)
			}
		}
	}
}
