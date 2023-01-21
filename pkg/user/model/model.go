package model

const (
	CUSTOMER_ROLE = 1
	COURIER_ROLE  = 2
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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
	Password string `json:"-"`
	Role     int    `json:"role"`
}

type MeResponse = User

type RegisterResponse = AuthTokens

type LoginResponse = AuthTokens

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
