package integration

import (
	"bira.io/template/dto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	Describe("Register", func() {
		It("Should successfully register user", func() {
			var body = dto.RegisterRequest{
				Username: "birareg",
				Password: "reg",
				Name:     "bira",
				Email:    "test@teste.com",
			}

			var response dto.RegisterResponse

			status, _ := Request(RequestObject{
				Method:      "POST",
				Path:        "/auth/register",
				RequestBody: body,
			}, &response)

			Expect(status).To(Equal(200))
			Expect(response.Message).To(Equal("Registered bira"))
		})
	})

	Describe("Login", func() {
		It("Should login", func() {
			var regBody = dto.RegisterRequest{
				Username: "biralogin",
				Password: "login",
				Name:     "bira",
				Email:    "test@teste.com",
			}

			status, _ := Request(RequestObject{
				Method:      "POST",
				Path:        "/auth/register",
				RequestBody: regBody,
			}, nil)

			Expect(status).To(Equal(200))

			var body = dto.LoginRequest{
				Username: "biralogin",
				Password: "login",
			}

			var response dto.LoginResponse

			status, _ = Request(RequestObject{
				Method:      "POST",
				Path:        "/auth/login",
				RequestBody: body,
			}, &response)

			Expect(status).To(Equal(200))
			Expect(len(response.Token)).ToNot(Equal(0))
		})
	})
})
