package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"kfamily/internal/model"
	"kfamily/internal/util"

	"github.com/google/uuid"
)

// This logic is repetitive, we already have user service,
// but we can't refer here, due to circular reference
// so have the same logic here

// GetUserByEmail matches case-insensitively: emails are stored lower-case
// (the "user" table has CHECK (email = lower(email))).
func (mw *Middleware) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	email = strings.ToLower(email)
	var dest model.User
	query := `SELECT * FROM "user" WHERE email = $1`
	if err := mw.db.GetContext(ctx, &dest, query, email); err != nil {
		return nil, fmt.Errorf("unable to get user by email. Err: %w", err)
	}

	return &dest, nil
}

// GetUserBySub looks up a user by their Logto subject (sub) claim, which is the
// stable, immutable identifier for a user - unlike email, which a user can change.
func (mw *Middleware) GetUserBySub(ctx context.Context, sub string) (*model.User, error) {
	var dest model.User
	query := `SELECT * FROM "user" WHERE sub = $1`
	err := mw.db.GetContext(ctx, &dest, query, sub)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("unable to get user by sub. Err: %w", err)
	}

	return &dest, nil
}

func (mw *Middleware) CreateUser(
	ctx context.Context,
	sub string,
	email string,
	name *string,
) (*model.User, error) {
	id := uuid.New()
	createdAt := time.Now()
	email = strings.ToLower(email)

	param := map[string]any{
		"id":         id,
		"sub":        sub,
		"email":      email,
		"name":       name,
		"created_at": createdAt,
		"created_by": id,
	}

	// "user" is a reserved word in Postgres and must be double-quoted in the
	// generated SQL (see internal/services/user/service.go's tableName).
	if err := util.Insert(ctx, mw.db, `"user"`, id, param); err != nil {
		return nil, fmt.Errorf("unable to create new user. err: %w", err)
	}

	return &model.User{
		ID:        id,
		Sub:       sub,
		Email:     email,
		Name:      name,
		IsAdmin:   false,
		CreatedAt: createdAt,
		CreatedBy: id,
	}, nil
}

// UpdateUserSub reassigns an existing local user's Logto sub. Logto can issue
// a different sub for the same email - e.g. signing in via a different
// method, or from a new browser/device before account linking has run - so a
// fresh sign-in can't always be matched to a local user by sub alone. See
// AuthMiddleware, which falls back to looking the user up by email and
// repoints their sub here instead of creating a duplicate (email is UNIQUE).
func (mw *Middleware) UpdateUserSub(ctx context.Context, existing *model.User, sub string) (*model.User, error) {
	updatedAt := time.Now()

	param := map[string]any{
		"id":         existing.ID,
		"sub":        sub,
		"updated_at": updatedAt,
		"updated_by": existing.ID,
	}

	// "user" is a reserved word in Postgres and must be double-quoted in the
	// generated SQL (see internal/services/user/service.go's tableName).
	if err := util.UpdateByID(ctx, mw.db, `"user"`, existing.ID, param); err != nil {
		return nil, fmt.Errorf("unable to update user sub. err: %w", err)
	}

	existing.Sub = sub
	existing.UpdatedAt = &updatedAt
	existing.UpdatedBy = &existing.ID
	return existing, nil
}
