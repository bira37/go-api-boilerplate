package unit

import (
	"time"

	repositoryContract "bira.io/template/contract/repository"
	serviceContract "bira.io/template/contract/service"
	"bira.io/template/dto"
	"bira.io/template/infra"
	"bira.io/template/model"
	"bira.io/template/service"
	repositoryMock "bira.io/template/test/mock/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Service", func() {
	var (
		userRepository repositoryContract.UserRepository
		userService    serviceContract.UserService
		ctrl           *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	Describe("GetLoggedUser", func() {
		Describe("With user in database", func() {
			BeforeEach(func() {
				userRepositoryMock := repositoryMock.NewMockUserRepository(ctrl)
				id, _ := uuid.NewUUID()
				userRepositoryMock.EXPECT().FindUserByUsername("birateste").Return(model.User{
					Username:     "birateste",
					Name:         "bira",
					Email:        "birateste@teste.com",
					Id:           id,
					PasswordHash: "hash",
					UpdatedAt:    time.Now().UTC(),
					CreatedAt:    time.Now().UTC(),
				}, nil)

				userRepositoryMock.EXPECT().FindUserByUsername("bira").Return(model.User{}, infra.NewSqlDbErrNotFound("User not found."))
				userRepository = userRepositoryMock

				userService = service.NewUserService(userRepository)
			})

			It("Should return the user if he is in database", func() {
				user, err := userService.GetLoggedUser(dto.GetLoggedUserRequest{Username: "birateste"})
				Expect(err).To(BeNil())
				Expect(user).NotTo(BeNil())
				Expect(user.Username).To(Equal("birateste"))
			})

			It("Should not return the user if he is not in the database", func() {
				user, err := userService.GetLoggedUser(dto.GetLoggedUserRequest{Username: "bira"})
				Expect(err).ToNot(BeNil())
				Expect(user).To(Equal(dto.GetLoggedUserResponse{}))
				Expect(err).To(BeAssignableToTypeOf(service.NewHttpErrBadRequest("")))
				Expect(err.Error()).To(Equal("User not found."))
			})
		})
	})
})
