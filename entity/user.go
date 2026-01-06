package entity

type User struct {
	ID          uint
	Name        string
	PhoneNumber string
	// password always keep hashed password.
	Password string
	Role     Role
}
