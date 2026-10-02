package auth

import (
	"kfamily/internal/telemetry"
	"kfamily/internal/util"
	"net/http"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/labstack/echo/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Signin exchanges an OAuth authorization code (obtained by the frontend
// redirecting the user to Casdoor's sign-in page) for a Casdoor access
// token. The exchange needs the application's client secret, so it has to
// happen here rather than in the browser.
//
// @Summary      Sign in with a Casdoor authorization code
// @Description  Rate limited per client IP.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body  auth.SigninInputDto  true  "Code and state from Casdoor's redirect"
// @Success      200  {object}  util.HttpResult{result=auth.TokenDto}
// @Failure      400,401,429  {object}  util.HttpResult
// @Router       /api/v1/auth/signin [post]
func (h *Handler) Signin(c *echo.Context) error {
	var input SigninInputDto
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, util.HttpError(err, "Unable to parse input values"))
	}

	if input.Code == "" {
		return c.JSON(http.StatusBadRequest, util.HttpErrorMessage("code is required"))
	}

	// The SDK builds its own request - this client puts it in our trace.
	httpClient := telemetry.HTTPClient(c.Request().Context())
	token, err := casdoorsdk.GetOAuthToken(input.Code, input.State, casdoorsdk.WithHTTPClient(httpClient))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, util.HttpError(err, "Unable to sign in with Casdoor"))
	}

	expiresIn := 0
	if !token.Expiry.IsZero() {
		expiresIn = int(time.Until(token.Expiry).Seconds())
	}

	return c.JSON(http.StatusOK, util.HttpData(TokenDto{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    expiresIn,
	}))
}
