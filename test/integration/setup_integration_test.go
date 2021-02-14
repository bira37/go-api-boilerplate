package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"bira.io/template/infra"
	"bira.io/template/server"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Tests")
}

var _ = Describe("Setup", func() {
	BeforeSuite(func() {
		sqldb := infra.NewSqlDb(infra.Config.SqlDbName)

		_ = sqldb.Execute(func(db *sqlx.DB) error {
			db.MustExec(`CREATE DATABASE IF NOT EXISTS testdb;`)
			return nil
		})
	})

	AfterSuite(func() {
		sqldb := infra.NewSqlDb(infra.Config.SqlDbName)

		_ = sqldb.Execute(func(db *sqlx.DB) error {
			db.MustExec(`DROP DATABASE IF EXISTS testdb;`)
			return nil
		})
	})
})

var once sync.Once

var (
	TestServer *httptest.Server
)

func GetServer() *httptest.Server {
	once.Do(func() {
		TestServer = httptest.NewServer(server.SetupServer("file://../../infra/migrations", "testdb"))
	})
	return TestServer
}

type RequestObject struct {
	Method      string
	Path        string
	RequestBody interface{}
	Headers     map[string]string
}

type ErrorBody struct {
	Code    string
	Message string
}

func Request(requestObject RequestObject, responseBody interface{}) (int, ErrorBody) {
	ts := GetServer()
	var body *bytes.Buffer = nil

	if requestObject.RequestBody != nil {
		jsonByte, err := json.Marshal(requestObject.RequestBody)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(jsonByte)
	}

	client := &http.Client{}

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(requestObject.Method, ts.URL+requestObject.Path, body)
	} else {
		req, err = http.NewRequest(requestObject.Method, ts.URL+requestObject.Path, nil)
	}

	if err != nil {
		panic(err)
	}

	if requestObject.Headers != nil {
		for k, v := range requestObject.Headers {
			req.Header.Set(k, v)
		}
	}

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	byteResponse, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	if responseBody != nil {
		err = json.Unmarshal(byteResponse, responseBody)

		if err != nil {
			panic(err)
		}
	}

	errorBody := ErrorBody{}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		err = json.Unmarshal(byteResponse, &errorBody)
	} else if responseBody != nil {
		err = json.Unmarshal(byteResponse, responseBody)
	}

	if err != nil {
		panic(err)
	}

	return res.StatusCode, errorBody
}
