package util

import (
	"context"
	"fmt"

	"kfamily/internal/dto"

	"github.com/jmoiron/sqlx"
)

// AdvancedQuery lists rows from tableName matching a dynamic WHERE clause
// built by the frontend's react-querybuilder (see AdvancedQueryBuilder.tsx):
// a parameterized SQL fragment referencing quoted column names and `:pN`
// placeholders, plus the parameter values in placeholder order. It's the
// shared implementation behind every table's "/query" endpoint (see
// kommune.Service.Query) - only the table name, select query and row/Dto
// types differ between tables.
//
// The fragment is validated against the table's real columns (via
// ValidateWhereCondition) and qualified with tableAlias (via
// QualifyWhereColumns) before use, since selectQuery may join in other
// tables (e.g. "user" for creator/updater/deleter emails) that share column
// names with the primary table.
//
// selectQuery must read from tableName aliased as tableAlias and scan into
// TRow; toDto converts each scanned row to its public Dto shape.
//
// Unlike List, this does NOT hide soft-deleted rows and ignores
// filter.IncludeDeletedRecords: the query builder exposes deleted_at /
// deleted_by like any other column, so callers filter on them explicitly
// (e.g. `"deleted_at" IS NULL`) when they want only live rows.
func AdvancedQuery[TRow any, TDto any](
	ctx context.Context,
	db *sqlx.DB,
	tableName string,
	tableAlias string,
	selectQuery string,
	filter dto.Filters,
	whereCondition string,
	whereConditionParams []any,
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

	columns, err := TableColumns(ctx, db, tableName)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s table metadata. err: %w", tableName, err)
	}
	allowedColumns := make(map[string]bool, len(columns))
	for _, c := range columns {
		allowedColumns[c.Name] = true
	}

	if err := ValidateWhereCondition(whereCondition, allowedColumns); err != nil {
		return nil, err
	}
	qualifiedWhereCondition := QualifyWhereColumns(whereCondition, tableAlias)

	params := map[string]any{
		"limit":  filter.ResultsPerPage,
		"offset": filter.PageIndex * filter.ResultsPerPage,
	}
	for i, v := range whereConditionParams {
		params[fmt.Sprintf("p%d", i+1)] = v
	}

	//total count
	var total int
	countQuery, countArgs, err := db.BindNamed(fmt.Sprintf(`SELECT COUNT(*) FROM %s %s WHERE `, tableName, tableAlias)+qualifiedWhereCondition, params)
	if err != nil {
		return nil, fmt.Errorf("unable to bind %s query count. err: %w", tableName, err)
	}
	if err := db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, fmt.Errorf("unable to run %s query count. err: %w", tableName, err)
	}

	//page of rows
	// id is a tie-breaker so OFFSET paging is stable - see util.List.
	listQuery, listArgs, err := db.BindNamed(fmt.Sprintf(`%s WHERE %s ORDER BY %s.%s %s, %s.id %s LIMIT :limit OFFSET :offset`,
		selectQuery, qualifiedWhereCondition, tableAlias, orderColumn, orderDirection, tableAlias, orderDirection), params)
	if err != nil {
		return nil, fmt.Errorf("unable to bind %s query list. err: %w", tableName, err)
	}

	rows := []TRow{}
	if err := db.SelectContext(ctx, &rows, listQuery, listArgs...); err != nil {
		return nil, fmt.Errorf("unable to run %s query. err: %w", tableName, err)
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
