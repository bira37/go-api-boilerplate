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

var _ = Describe("Auth Service", func() {
	var (
		userRepository repositoryContract.UserRepository
		authService    serviceContract.AuthService
		ctrl           *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	Describe("Login", func() {
		Describe("With user in database", func() {
			BeforeEach(func() {
				userRepositoryMock := repositoryMock.NewMockUserRepository(ctrl)
				id, _ := uuid.NewUUID()
				userRepositoryMock.EXPECT().FindUserByUsername("birateste").Return(model.User{
					Username:     "birateste",
					Name:         "bira",
					Email:        "birateste@teste.com",
					Id:           id,
					PasswordHash: "$2y$10$ksczaxsb4JWTOrwQ9lh3c.y6Rt85d1OD3MwaXjN1J/9cPTNMcUtxy",
					UpdatedAt:    time.Now().UTC(),
					CreatedAt:    time.Now().UTC(),
				}, nil)

				userRepositoryMock.EXPECT().FindUserByUsername("bira").Return(model.User{}, infra.NewSqlDbErrNotFound("User not found."))
				userRepository = userRepositoryMock

				authService = service.NewAuthService(userRepository)
			})

			It("Should login with valid credentials", func() {
				res, err := authService.Login(dto.LoginRequest{Username: "birateste", Password: "hash"})
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.Token).NotTo(BeNil())
			})

			It("Should not login with invalid password", func() {
				res, err := authService.Login(dto.LoginRequest{Username: "birateste", Password: "hasherrado"})
				Expect(err).ToNot(BeNil())
				Expect(res).To(Equal(dto.LoginResponse{}))
				Expect(err).To(BeAssignableToTypeOf(service.NewHttpError("", "", 0)))
				Expect(err.Error()).To(Equal("Invalid login or password."))
			})

			It("Should not login with invalid username", func() {
				res, err := authService.Login(dto.LoginRequest{Username: "bira", Password: "hash"})
				Expect(err).ToNot(BeNil())
				Expect(res).To(Equal(dto.LoginResponse{}))
				Expect(err).To(BeAssignableToTypeOf(service.NewHttpError("", "", 0)))
				Expect(err.Error()).To(Equal("User not found."))
			})
		})

		Describe("With malfunctioned database", func() {
			BeforeEach(func() {
				userRepositoryMock := repositoryMock.NewMockUserRepository(ctrl)
				userRepositoryMock.EXPECT().FindUserByUsername("bira").Return(model.User{}, infra.NewSqlDbErrInternal("Internal error."))
				userRepository = userRepositoryMock

				authService = service.NewAuthService(userRepository)
			})

			It("Should return internal error", func() {
				_, err := authService.Login(dto.LoginRequest{Username: "bira", Password: "hash"})
				Expect(err).ToNot(BeNil())
				Expect(err).To(BeAssignableToTypeOf(service.NewHttpError("", "", 0)))
				Expect(err.Error()).To(Equal("Internal error."))
			})
		})
	})

	Describe("Register", func() {
		Describe("With user in database", func() {
			BeforeEach(func() {
				userRepositoryMock := repositoryMock.NewMockUserRepository(ctrl)
				id, _ := uuid.NewUUID()
				id2, _ := uuid.NewUUID()
				userRepositoryMock.EXPECT().FindUserByUsername("birateste").Return(model.User{
					Username:     "birateste",
					Name:         "bira",
					Email:        "birateste@teste.com",
					Id:           id,
					PasswordHash: "$2y$10$ksczaxsb4JWTOrwQ9lh3c.y6Rt85d1OD3MwaXjN1J/9cPTNMcUtxy",
					UpdatedAt:    time.Now().UTC(),
					CreatedAt:    time.Now().UTC(),
				}, nil)

				userRepositoryMock.EXPECT().FindUserByUsername("bira").Return(model.User{}, infra.NewSqlDbErrNotFound("User not found."))
				userRepositoryMock.EXPECT().InsertUser(gomock.Any()).Return(model.User{Username: "bira", PasswordHash: "hash", Name: "bira", Email: "test@teste.com", Id: id2, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}, nil)

				userRepository = userRepositoryMock

				authService = service.NewAuthService(userRepository)
			})

			It("should register successfully a new user", func() {
				res, err := authService.Register(dto.RegisterRequest{Username: "bira", Password: "hash", Name: "bira", Email: "test@teste.com"})
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.Message).To(Equal("Registered bira"))
			})

			It("Should not register user if already exists one with the same username", func() {
				_, err := authService.Register(dto.RegisterRequest{Username: "birateste", Password: "hash", Name: "bira", Email: "test@teste.com"})
				Expect(err).NotTo(BeNil())
				Expect(err).To(BeAssignableToTypeOf(service.NewHttpError("", "", 0)))
				Expect(err.Error()).To(Equal("An user with the same username already exists."))
			})
		})
	})
})
