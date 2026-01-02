package uservalidator

import (
	"errors"
	"fmt"
	"game_app/dto"
	"game_app/pkg/errmsg"
	"game_app/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateLoginRequest(req dto.LoginRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required),

		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(PhoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.doesPhoneNumberExist)),
	); err != nil {
		filedErrors := make(map[string]string)

		var errV validation.Errors
		if ok := errors.As(err, &errV); ok {
			for key, value := range errV {
				if value != nil {
					filedErrors[key] = value.Error()
				}
			}
		}

		return richerror.New(op).
			SetMessage(errmsg.ErrorMsgInvalidInput).
			SetKind(richerror.KindInvalid).
			SetMeta(map[string]interface{}{"request": req}).
			SetWrappedError(err), filedErrors
	}

	return nil, nil
}

func (v Validator) doesPhoneNumberExist(value interface{}) error {
	phoneNumber := value.(string)
	_, err := v.repository.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf(errmsg.ErrorMsgNotFound)
	}

	return nil
}
