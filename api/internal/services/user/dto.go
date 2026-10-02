package user

import "kfamily/internal/dto"

type (
	// UpdateInputDto is the only mutable shape for a user - everything else
	// (sub, email) is fixed at creation time by the auth middleware.
	// IsAdmin is a pointer so that omitting it leaves the flag unchanged,
	// rather than silently demoting the user.
	UpdateInputDto struct {
		Name    *string `json:"name"`
		IsAdmin *bool   `json:"isAdmin"`
	}

	// Dto is the API response shape for a user.
	Dto struct {
		Id      string  `json:"id"`
		Sub     string  `json:"sub"`
		Email   string  `json:"email"`
		Name    *string `json:"name"`
		IsAdmin bool    `json:"isAdmin"`
		// A service account (batch job, system) - never signs in.
		IsService bool `json:"isService"`
		dto.AuditDto
	}
)

// Validate trims the input and checks it fits the user table's columns.
// name is an optional varchar(100) - a blank name is stored as NULL.
func (i *UpdateInputDto) Validate() error {
	if i.Name == nil {
		return nil
	}

	name, err := dto.TrimmedText("name", *i.Name, false, 100)
	if err != nil {
		return err
	}
	if name == "" {
		i.Name = nil
	} else {
		i.Name = &name
	}

	return nil
}
