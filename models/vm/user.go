package vm

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" form:"pn"`
}

type LoginResponse struct {
	UserID      uint64 `json:"user_id"`
	AccessToken string `json:"access_token"`
}
