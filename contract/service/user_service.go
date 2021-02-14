package contract

import "bira.io/template/dto"

type UserService interface {
	GetLoggedUser(dto.GetLoggedUserRequest) (dto.GetLoggedUserResponse, error)
}
