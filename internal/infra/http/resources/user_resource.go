package resources

import "go-rest-api/internal/domain"

type UserDto struct {
	Id     uint64  `json:"id"`
	Name   string  `json:"username"`
	Email  string  `json:"email"`
	Avatar *string `json:"avatar"`
}

func (u UserDto) DomainToDto(user domain.User) UserDto {
	return UserDto{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Avatar: user.Avatar,
	}
}
