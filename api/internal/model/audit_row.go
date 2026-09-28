package model

import (
	"time"

	"kfamily/internal/dto"

	"github.com/google/uuid"
)

// AuditRow is the flat scan target for the created_by/updated_by/deleted_by
// trio every audited table joins in its creator/updater/deleter email for
// (see util.AuditSelectColumns/AuditSelectJoins). Embed it in a table's row
// struct so sqlx scans these columns automatically, and call ToDto to fill
// in a dto.AuditDto.
type AuditRow struct {
	CreatedAt      time.Time  `db:"created_at"`
	CreatedByID    uuid.UUID  `db:"created_by"`
	CreatedByEmail string     `db:"created_by_email"`
	UpdatedAt      *time.Time `db:"updated_at"`
	UpdatedByID    *uuid.UUID `db:"updated_by"`
	UpdatedByEmail *string    `db:"updated_by_email"`
	DeletedAt      *time.Time `db:"deleted_at"`
	DeletedByID    *uuid.UUID `db:"deleted_by"`
	DeletedByEmail *string    `db:"deleted_by_email"`
}

func (a AuditRow) ToDto() dto.AuditDto {
	d := dto.AuditDto{
		CreatedAt: a.CreatedAt,
		CreatedBy: dto.IdEmail{Id: a.CreatedByID.String(), Email: a.CreatedByEmail},
		UpdatedAt: a.UpdatedAt,
		DeletedAt: a.DeletedAt,
	}

	d.UpdatedBy = idEmail(a.UpdatedByID, a.UpdatedByEmail)
	d.DeletedBy = idEmail(a.DeletedByID, a.DeletedByEmail)

	return d
}

// idEmail builds an optional IdEmail from a nullable *_by column and its
// LEFT JOINed email, without dereferencing either blindly - a nil pointer
// here would otherwise panic the whole request.
func idEmail(id *uuid.UUID, email *string) *dto.IdEmail {
	if id == nil {
		return nil
	}

	result := &dto.IdEmail{Id: id.String()}
	if email != nil {
		result.Email = *email
	}
	return result
}
