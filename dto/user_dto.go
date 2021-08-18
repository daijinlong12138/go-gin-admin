package dto

import "go-gin-admin/model"

type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(users model.Users) UserDto {
	return UserDto{
		Name:      users.Name,
		Telephone: users.Telephone,
	}
}
