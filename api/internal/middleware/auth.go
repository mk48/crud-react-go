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

		sub := claims.User.Id
		if sub == "" {
			return c.JSON(http.StatusUnauthorized, util.HttpErrorMessage("Unauthorized. Token missing subject"))
		}

		ctx := c.Request().Context()

		dbUser, err := mw.GetUserBySub(ctx, sub)
		if err == sql.ErrNoRows {
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
				dbUser, err = mw.UpdateUserSub(ctx, existingUser, sub)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to update existing user's sub"))
				}
			} else if errors.Is(err, sql.ErrNoRows) {
				name := claims.User.DisplayName
				dbUser, err = mw.CreateUser(ctx, sub, claims.User.Email, &name)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to create the new user"))
				}
			} else {
				return c.JSON(http.StatusInternalServerError, util.HttpError(err, "Unable to look up user by email"))
			}
		} else if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusUnauthorized, util.HttpError(err, "Error to get user by sub"))
		}

		// Users are soft-deleted (see user.Service.Delete), so a deleted user
		// still resolves above - by sub, or by email reconciliation.
		if dbUser.DeletedAt != nil {
			return c.JSON(http.StatusForbidden, util.HttpErrorMessage("Forbidden. User account is deleted"))
		}

		// set user in context
		c.Set("user", dbUser)

		// Call the next handler in the chain
		return next(c)
	}
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
