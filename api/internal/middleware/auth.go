package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"kfamily/internal/util"
	"net/http"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/labstack/echo/v5"
)

func (mw *Middleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {

		authHeader := c.Request().Header.Values("Authorization")
		if len(authHeader) == 0 {
			return c.JSON(http.StatusUnauthorized, util.HttpErrorMessage("Unauthorized. Authorization header missing"))
		}

		var authorization = authHeader[0]
		var sessionToken = strings.TrimPrefix(authorization, "Bearer ")

		// Verify the token's signature against Casdoor's certificate. The
		// resulting claims already carry the full Casdoor user profile, so
		// there's no separate call to a management API needed to provision a
		// local user below.
		claims, err := casdoorsdk.ParseJwtToken(sessionToken)
		if err != nil {
			// Include the jwt error (e.g. "token is expired" vs
			// "crypto/rsa: verification error" for a certificate that doesn't
			// match the one the Casdoor application signs with) - otherwise
			// these cases are indistinguishable from the client.
			return c.JSON(http.StatusUnauthorized, util.HttpError(err, "Unauthorized. Can't verify token"))
		}

		client, err := mw.validateClaims(claims)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, util.HttpError(err, "Unauthorized. Token not valid for this application"))
		}

		// Which app this request comes from, for every operation it runs
		// (see util.RunOperation) - set before anything below can write,
		// e.g. provisioning a first-time user.
		ctx := util.WithRequestSource(c.Request().Context(), util.RequestSource{
			Client: client.Name,
			Info:   clientInfo(c),
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// A batch job's token (client-credentials grant) has no signed-in
		// user - it acts as its client's service account.
		if client.IsService() {
			serviceUser, err := mw.GetUserByID(ctx, *client.ServiceUserID)
			if err != nil {
				c.Logger().Error("unable to load service account", "client", client.Name, "err", err)
				return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to load the client's service account"))
			}
			if !serviceUser.IsService || serviceUser.DeletedAt != nil {
				return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. Client's service account is not active"))
			}

			c.Set("user", serviceUser)
			return next(c)
		}

		sub := claims.User.Id
		if sub == "" {
			return c.JSON(http.StatusUnauthorized, util.HttpErrorMessage("Unauthorized. Token missing subject"))
		}

		dbUser, err := mw.GetUserBySub(ctx, sub)
		if errors.Is(err, sql.ErrNoRows) {
			// First time we've seen this sub - provision a local user record
			// for them using the profile embedded in the token.
			if claims.User.Email == "" {
				return c.JSON(http.StatusUnauthorized, util.HttpErrorMessage("Unauthorized. Casdoor user has no email"))
			}

			// Casdoor can issue a different sub for the same email (e.g.
			// signing in via a different method, or from a new
			// browser/device before account linking has run). Reconcile
			// onto the existing local user by email instead of creating a
			// duplicate, which would violate the unique constraint on email.
			existingUser, err := mw.GetUserByEmail(ctx, claims.User.Email)
			if err == nil {
				// Deleted rows are read-only (see util.UpdateByID), so refuse
				// here rather than failing the sub update below.
				if existingUser.DeletedAt != nil {
					return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. User account is deleted"))
				}
				// A Casdoor user sharing a service account's email must not
				// take it over.
				if existingUser.IsService {
					return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. Service accounts can't sign in"))
				}
				dbUser, err = mw.UpdateUserSub(ctx, existingUser, sub)
				if err != nil {
					c.Logger().Error("unable to update existing user's sub", "err", err)
					return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to update existing user's sub"))
				}
			} else if errors.Is(err, sql.ErrNoRows) {
				// "user".name is NULL or non-blank (CHECK constraint).
				var name *string
				if displayName := strings.TrimSpace(claims.User.DisplayName); displayName != "" {
					name = &displayName
				}
				dbUser, err = mw.CreateUser(ctx, sub, claims.User.Email, name)
				if err != nil {
					// A concurrent first request for the same sub may have
					// created the user between our lookup and insert (the
					// insert then fails on the unique sub) - use theirs.
					if existing, lookupErr := mw.GetUserBySub(ctx, sub); lookupErr == nil {
						dbUser = existing
					} else {
						c.Logger().Error("unable to create the new user", "err", err)
						return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to create the new user"))
					}
				}
			} else {
				c.Logger().Error("unable to look up user by email", "err", err)
				return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to look up user by email"))
			}
		} else if err != nil {
			c.Logger().Error("unable to get user by sub", "err", err)
			return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to get user by sub"))
		}

		// Users are soft-deleted (see user.Service.Delete), so a deleted user
		// still resolves above - by sub, or by email reconciliation.
		if dbUser.DeletedAt != nil {
			return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. User account is deleted"))
		}
		if dbUser.IsService {
			return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. Service accounts can't sign in"))
		}

		// set user in context
		c.Set("user", dbUser)

		// Call the next handler in the chain
		return next(c)
	}
}

// validateClaims checks that a signature-verified token was issued by our
// Casdoor, as an access token, to one of the registered apps (see
// util.ClientRegistry) - and returns that app. ParseJwtToken only verifies
// the signature and expiry - and Casdoor applications commonly share the
// built-in signing certificate - so a validly signed token could otherwise
// belong to an unrelated application, or be a (long-lived) refresh token
// presented as an access token.
func (mw *Middleware) validateClaims(claims *casdoorsdk.Claims) (util.Client, error) {
	if claims.TokenType != "access-token" {
		return util.Client{}, fmt.Errorf("unexpected token type %q", claims.TokenType)
	}

	if strings.TrimSuffix(claims.Issuer, "/") != strings.TrimSuffix(mw.env.CasdoorEndpoint, "/") {
		return util.Client{}, fmt.Errorf("unexpected token issuer %q", claims.Issuer)
	}

	// azp ("authorized party") names the application the token was issued
	// to; fall back to the audience when Casdoor leaves it out.
	if claims.Azp != "" {
		if client, ok := mw.clients.Lookup(claims.Azp); ok {
			return client, nil
		}
		return util.Client{}, fmt.Errorf("token was issued to unregistered application %q", claims.Azp)
	}
	for _, aud := range claims.Audience {
		if client, ok := mw.clients.Lookup(aud); ok {
			return client, nil
		}
	}

	return util.Client{}, fmt.Errorf("token audience %v is not a registered application", claims.Audience)
}

// clientInfo is what the caller reports about itself (see
// util.RequestSource.Info) - capped in length, since it's caller-controlled
// and stored on every operation.
func clientInfo(c *echo.Context) map[string]any {
	info := map[string]any{
		"ip":        c.RealIP(),
		"userAgent": truncate(c.Request().UserAgent(), 256),
	}
	if v := c.Request().Header.Get(util.ClientVersionHeader); v != "" {
		info["version"] = truncate(v, 64)
	}
	if p := c.Request().Header.Get(util.ClientPlatformHeader); p != "" {
		info["platform"] = truncate(p, 32)
	}
	return info
}

func truncate(s string, max int) string {
	if r := []rune(s); len(r) > max {
		return string(r[:max])
	}
	return s
}

// AdminMiddleware restricts access to admin users. It must run after
// AuthMiddleware, which sets the "user" context value it relies on.
func (mw *Middleware) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		user := util.CurrentUser(c)
		if !user.IsAdmin {
			return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. Admin access required"))
		}

		return next(c)
	}
}
