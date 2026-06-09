package dto

type RegisterRequest struct {
	Body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`

		PubKey     []byte `json:"pubKey"`
		EncPrivKey []byte `json:"encPrivKey"`
	}
}

type RegisterResponse struct {
	Status int
}

type LoginRequest struct {
	Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}

// According to rfc6750 this is what should be included in response
type LoginResponse struct {
	Status int
	Body   LoginResponseBody
}

type LoginResponseBody struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
