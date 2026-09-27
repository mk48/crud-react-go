package sample

import "kfamily/internal/dto"

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
