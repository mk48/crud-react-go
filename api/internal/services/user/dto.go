package user

import "kfamily/internal/dto"

type (
	// UpdateInputDto is the only mutable shape for a user - everything else
	// (sub, email) is fixed at creation time by the auth middleware.
	UpdateInputDto struct {
		Name    *string `json:"name"`
		IsAdmin bool    `json:"isAdmin"`
	}

	// Dto is the API response shape for a user.
	Dto struct {
		Id      string  `json:"id"`
		Sub     string  `json:"sub"`
		Email   string  `json:"email"`
		Name    *string `json:"name"`
		IsAdmin bool    `json:"isAdmin"`
		dto.AuditDto
	}
)
