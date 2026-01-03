package uservalidator

import (
	"errors"
	"fmt"
	"game_app/param"
	"game_app/pkg/errmsg"
	"game_app/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),

		validation.Field(&req.Password,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#$%&*]{8,}$`)).
				Error(errmsg.ErrorMsgPhoneNumberIsNotValid)),

		validation.Field(&req.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(PhoneNumberRegex)),
			validation.By(v.checkPhoneNumberUniqueness)),
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

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)

	if isUnique, err := v.repository.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			return err
		}

		return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
	}

	return nil
}
