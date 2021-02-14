package integration

import (
	"bira.io/template/dto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	Describe("GetLoggedUser", func() {
		It("Should get me when logged", func() {
			var regBody = dto.RegisterRequest{
				Username: "biraget",
				Password: "getme",
				Name:     "bira",
				Email:    "test@teste.com",
			}

			status, _ := Request(RequestObject{
				Method:      "POST",
				Path:        "/auth/register",
				RequestBody: regBody,
			}, nil)

			Expect(status).To(Equal(200))

			var loginBody = dto.LoginRequest{
				Username: "biraget",
				Password: "getme",
			}

			var loginResponse dto.LoginResponse

			status, _ = Request(RequestObject{
				Method:      "POST",
				Path:        "/auth/login",
				RequestBody: loginBody,
			}, &loginResponse)

			Expect(status).To(Equal(200))

			var response dto.GetLoggedUserResponse

			status, _ = Request(RequestObject{
				Method: "GET",
				Path:   "/user/me",
				Headers: map[string]string{
					"X-Access-Token": loginResponse.Token,
				},
			}, &response)

			Expect(status).To(Equal(200))
		})
	})
})
