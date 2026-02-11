package authorizationservice

import (
	"fmt"
	"game_app/entity"
	"game_app/pkg/richerror"
)

type Repository interface {
	GetUserPermissionsTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) CheckAccess(userID uint, role entity.Role, permission ...entity.PermissionTitle) (bool, error) {
	const op = "authorizationservice.CheckAccess"
	permissionsTitle, err := s.repo.GetUserPermissionsTitles(userID, role)
	if err != nil {
		fmt.Println("opp: ", op)
		return false, richerror.New(op).SetWrappedError(err)
	}

	// check the access

	for _, pt := range permissionsTitle {
		fmt.Println("for1")
		for _, p := range permission {
			fmt.Println("for2")
			if p == pt {
				return true, nil
			}
		}
	}

	return false, nil
}
