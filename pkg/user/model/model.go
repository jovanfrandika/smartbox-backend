package model

type MeInput struct {
	ID string `json:"id"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type MeResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  int    `json:"role"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
