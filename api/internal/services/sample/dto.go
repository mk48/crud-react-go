package sample

import (
	"strings"

	"kfamily/internal/dto"
)

type (
	CreateInputDto struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	UpdateInputDto struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	// Dto is the API response shape for a sample item.
	Dto struct {
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		dto.AuditDto
	}
)

// Validate trims the input and checks it fits sample_items' columns.
func (i *CreateInputDto) Validate() error {
	return validate(&i.Name, &i.Description)
}

// Validate trims the input and checks it fits sample_items' columns.
func (i *UpdateInputDto) Validate() error {
	return validate(&i.Name, &i.Description)
}

func validate(name *string, description **string) error {
	var err error
	if *name, err = dto.TrimmedText("name", *name, true, 200); err != nil {
		return err
	}

	// description is an optional TEXT column - store blank as NULL.
	if *description != nil {
		trimmed := strings.TrimSpace(**description)
		if trimmed == "" {
			*description = nil
		} else {
			*description = &trimmed
		}
	}

	return nil
}
