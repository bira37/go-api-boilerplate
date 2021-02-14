package contract

import "bira.io/template/dto"

type AuthService interface {
	Login(dto.LoginRequest) (dto.LoginResponse, error)
	Register(dto.RegisterRequest) (dto.RegisterResponse, error)
}
