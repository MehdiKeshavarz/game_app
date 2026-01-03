package param

type GetProfileRequest struct {
	UserID uint `json:"user_id"`
}
type GetProfileResponse struct {
	Name string `json:"name"`
}
