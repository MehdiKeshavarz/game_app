package userservice

import (
	"game_app/dto"
	"game_app/pkg/richerror"
)

func (s Service) Profile(req dto.GetProfileRequest) (dto.GetProfileResponse, error) {
	const op = "userservice.Profile"
	// I don't expect the repository call return "record not found " error,
	// because I assume the interactor input is sanitized.

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.GetProfileResponse{}, richerror.New(op).
			SetWrappedError(err).
			SetMeta(map[string]interface{}{"req": req})
	}

	return dto.GetProfileResponse{Name: user.Name}, nil
}
