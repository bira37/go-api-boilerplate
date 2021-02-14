package service

import (
	"time"

	"bira.io/template/infra"
	"bira.io/template/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	repositoryContract "bira.io/template/contract/repository"
	serviceContract "bira.io/template/contract/service"
	"bira.io/template/dto"
)

type authService struct {
	userRepository repositoryContract.UserRepository
}

func NewAuthService(ur repositoryContract.UserRepository) serviceContract.AuthService {
	return &authService{
		userRepository: ur,
	}
}

func (s *authService) Login(loginRequest dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepository.FindUserByUsername(loginRequest.Username)

	if err != nil {
		dberr, _ := err.(*infra.SqlDbError)
		if dberr.Code == infra.ErrDbNotFound {
			return dto.LoginResponse{}, NewHttpErrNotFound("User not found.")
		}
		return dto.LoginResponse{}, NewHttpErrInternalServer("Internal error.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password))

	if err != nil {
		return dto.LoginResponse{}, NewHttpErrBadRequest("Invalid login or password.")
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims = jwt.MapClaims{
		"exp": time.Now().UTC().Add(time.Hour * 72).Unix(),
		"iat": time.Now().UTC(),
		"sub": loginRequest.Username,
	}

	tokenString, err := token.SignedString([]byte(infra.Config.JwtSigningString))

	if err != nil {
		return dto.LoginResponse{}, NewHttpErrInternalServer("Internal error.")
	}

	return dto.LoginResponse{
		Message: "Hello, " + user.Name,
		Token:   tokenString,
	}, nil
}

func (s *authService) Register(registerRequest dto.RegisterRequest) (dto.RegisterResponse, error) {
	_, err := s.userRepository.FindUserByUsername(registerRequest.Username)

	if err == nil {
		return dto.RegisterResponse{}, NewHttpErrBadRequest("An user with the same username already exists.")
	}

	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		return dto.RegisterResponse{}, err
	}

	hashedPassword := string(hashedPasswordByte)

	user, err := s.userRepository.InsertUser(model.UserCreate{Name: registerRequest.Name, Username: registerRequest.Username, PasswordHash: hashedPassword, Email: registerRequest.Email})

	if err != nil {
		return dto.RegisterResponse{}, NewHttpErrInternalServer("Internal error.")
	}

	return dto.RegisterResponse{
		Message: "Registered " + user.Name,
	}, nil
}
