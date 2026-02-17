package matchingvalidator

import (
	"errors"
	"fmt"
	"game_app/entity"
	"game_app/param"
	"game_app/pkg/errmsg"
	"game_app/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddToWaitingList(req param.AddToWaitingListRequest) (error, map[string]string) {
	const op = "matchingvalidator.ValidateAddToWaitingList"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.Category,
			validation.Required,
			validation.By(v.isCategoryValid)),
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

func (v Validator) isCategoryValid(value interface{}) error {
	category := value.(entity.Category)

	if !category.IsValid() {
		return fmt.Errorf(errmsg.ErrorMsgCategoryIsNotValid)
	}

	return nil
}
