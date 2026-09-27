package auth

type (
	SigninInputDto struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	// TokenDto is what the frontend stores and sends back as a Bearer token.
	// AccessToken is a Casdoor-issued JWT - internal/middleware.AuthMiddleware
	// verifies it directly against Casdoor's certificate, so the API never
	// needs its own session/token format.
	TokenDto struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		TokenType    string `json:"tokenType"`
		ExpiresIn    int    `json:"expiresIn"`
	}
)
