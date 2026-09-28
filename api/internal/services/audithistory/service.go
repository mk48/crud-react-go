package audithistory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

//------------------------------------------------------------------------------

// List returns every audit_history record for sourceId — a source table's
// primary key — most recent first.
func (s *Service) List(ctx context.Context, sourceId uuid.UUID) ([]Dto, error) {
	rows := []row{}
	query := selectQuery + ` WHERE source_id = $1 ORDER BY created_at DESC, id DESC`
	if err := s.db.SelectContext(ctx, &rows, query, sourceId); err != nil {
		return nil, fmt.Errorf("unable to load audit history. err: %w", err)
	}

	items := make([]Dto, len(rows))
	for i, r := range rows {
		items[i] = r.toDto()
	}

	return items, nil
}
