package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

const (
	UserRoelStr  = "user"
	AdminRoleStr = "admin"
)

func (r Role) String() string {
	switch r {
	case UserRole:
		return UserRoelStr
	case AdminRole:
		return AdminRoleStr
	}

	return ""
}

func MapToRole(roleStr string) Role {
	switch roleStr {
	case UserRoelStr:
		return UserRole
	case AdminRoleStr:
		return AdminRole

	}

	return Role(0)
}
