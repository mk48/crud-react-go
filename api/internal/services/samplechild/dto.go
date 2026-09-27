package samplechild

import "kfamily/internal/dto"

type (
	CreateInputDto struct {
		SampleItemId string `json:"sampleItemId"`
		Name         string `json:"name"`
	}

	UpdateInputDto struct {
		SampleItemId string `json:"sampleItemId"`
		Name         string `json:"name"`
	}

	// Dto is the API response shape for a sample child item. The parent
	// sample item is returned as {id, name} (joined in by selectQuery) so the
	// frontend can show/pre-select it without a second request.
	Dto struct {
		Id         string     `json:"id"`
		SampleItem dto.IdName `json:"sampleItem"`
		Name       string     `json:"name"`
		dto.AuditDto
	}
)
