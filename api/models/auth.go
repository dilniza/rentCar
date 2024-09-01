package models

type CustomerLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CustomerLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthInfo struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type ChangePassword struct {
	Login       string `json:"login"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CustomerRegisterRequest struct {
	Mail string `json:"mail"`
}

type CustomerRegisterConfirm struct {
	Mail     string         `json:"mail"`
	Otp      string         `json:"otp"`
	Customer CreateCustomer `json:"customer"`
}
