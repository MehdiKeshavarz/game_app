package uservalidator

import "game_app/entity"

const (
	PhoneNumberRegex = "^(0|0098|\\+98)9(0[1-5]|[1 3]\\d|2[0-2]|98)\\d{7}$"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repository Repository
}

func New(repository Repository) Validator {
	return Validator{repository: repository}
}
