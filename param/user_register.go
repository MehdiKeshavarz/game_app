package param

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type RegisterResponse struct {
	User struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user"`
}
