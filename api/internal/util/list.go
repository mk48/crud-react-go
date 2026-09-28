package util

import (
	"context"
	"fmt"
	"strings"

	"kfamily/internal/dto"

	"github.com/jmoiron/sqlx"
)

// List returns a page of rows from tableName ordered by filter.SortColumn
// (already validated against the table's allowed sort keys by
// PaginationFilter - see each table's init.go), restricted to non-deleted
// rows unless filter.IncludeDeletedRecords is set, and optionally matching a
// free-text search across searchColumns (pass nil/empty if the table has no
// search box). It's the shared implementation behind every table's plain
// list endpoint (see kommune.Service.List) - only the table name, select
// query, searchable columns and row/Dto types differ between tables.
//
// selectQuery must read from tableName aliased as tableAlias and scan into
// TRow, and TRow's table must have a deleted_at column (every audited table
// does - see util.RecordAudit); toDto converts each scanned row to its
// public Dto shape.
func List[TRow any, TDto any](
	ctx context.Context,
	db *sqlx.DB,
	tableName string,
	tableAlias string,
	selectQuery string,
	searchColumns []string,
	filter dto.Filters,
	toDto func(TRow) TDto,
) (*dto.PaginationResponse[TDto], error) {
	orderColumn := "created_at"
	if filter.SortColumn != "" {
		orderColumn = filter.SortColumn
	}
	orderDirection := "ASC"
	if filter.SortDirection == "desc" {
		orderDirection = "DESC"
	}

	// Defaults to TRUE (matches every row) when the table has no searchable
	// columns or no search term, so the rest of the WHERE clause is
	// unaffected. A "%%" ILIKE would otherwise exclude rows with a NULL in
	// the searched column, since NULL ILIKE '%%' evaluates to NULL, not TRUE.
	searchClause := "TRUE"
	if len(searchColumns) > 0 && filter.Search != "" {
		parts := make([]string, len(searchColumns))
		for i, c := range searchColumns {
			parts[i] = fmt.Sprintf(`%s."%s" ILIKE :search`, tableAlias, c)
		}
		searchClause = "(" + strings.Join(parts, " OR ") + ")"
	}

	params := map[string]any{
		"search": "%" + filter.Search + "%",
		"limit":  filter.ResultsPerPage,
		"offset": filter.PageIndex * filter.ResultsPerPage,
	}

	// Built statically (not "(:include_deleted OR ...)") so the planner can
	// use the live-rows partial indexes (WHERE deleted_at IS NULL).
	whereClause := searchClause
	if !filter.IncludeDeletedRecords {
		whereClause = fmt.Sprintf(`%s.deleted_at IS NULL AND %s`, tableAlias, searchClause)
	}

	//total count
	var total int
	countQuery, countArgs, err := db.BindNamed(fmt.Sprintf(`SELECT COUNT(*) FROM %s %s WHERE `, tableName, tableAlias)+whereClause, params)
	if err != nil {
		return nil, fmt.Errorf("unable to bind %s list count. err: %w", tableName, err)
	}
	if err := db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, fmt.Errorf("unable to count %s list. err: %w", tableName, err)
	}

	//page of rows
	// id is a tie-breaker: without a unique last sort key, rows sharing the
	// sort value can come back in a different order per query, so OFFSET
	// paging would repeat or skip them across pages. Same direction as the
	// sort, so one (column, id) index serves both ASC and DESC.
	listQuery, listArgs, err := db.BindNamed(fmt.Sprintf(`%s WHERE %s ORDER BY %s.%s %s, %s.id %s LIMIT :limit OFFSET :offset`,
		selectQuery, whereClause, tableAlias, orderColumn, orderDirection, tableAlias, orderDirection), params)
	if err != nil {
		return nil, fmt.Errorf("unable to bind %s list query. err: %w", tableName, err)
	}

	rows := []TRow{}
	//fmt.Println("List query:", debugInterpolateSQL(listQuery, listArgs))
	if err := db.SelectContext(ctx, &rows, listQuery, listArgs...); err != nil {
		return nil, fmt.Errorf("unable to load %s list. err: %w", tableName, err)
	}

	items := make([]TDto, len(rows))
	for i, r := range rows {
		items[i] = toDto(r)
	}

	return &dto.PaginationResponse[TDto]{
		Items: items,
		Pagination: dto.PaginationInfo{
			PageIndex:      filter.PageIndex,
			ResultsPerPage: filter.ResultsPerPage,
			TotalResults:   total,
		},
	}, nil
}
