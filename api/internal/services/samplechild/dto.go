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

// Validate trims the input and checks it fits sample_child_items' columns.
// sampleItemId is checked separately by the service (it needs the database).
func (i *CreateInputDto) Validate() error {
	var err error
	i.Name, err = dto.TrimmedText("name", i.Name, true, 200)
	return err
}

// Validate trims the input and checks it fits sample_child_items' columns.
// sampleItemId is checked separately by the service (it needs the database).
func (i *UpdateInputDto) Validate() error {
	var err error
	i.Name, err = dto.TrimmedText("name", i.Name, true, 200)
	return err
}
