package dto

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type (
	Id struct {
		Id string `json:"id"`
	}

	IdEmail struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}

	IdName struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	AuditDto struct {
		CreatedAt time.Time  `json:"createdAt"`
		CreatedBy IdEmail    `json:"createdBy"`
		UpdatedAt *time.Time `json:"updatedAt"`
		UpdatedBy *IdEmail   `json:"updatedBy"`
		DeletedAt *time.Time `json:"deletedAt"`
		DeletedBy *IdEmail   `json:"deletedBy"`
	}

	PaginationInfo struct {
		PageIndex      int `json:"pageIndex"`
		ResultsPerPage int `json:"resultsPerPage"`
		TotalResults   int `json:"totalResults"`
	}

	PaginationResponse[T any] struct {
		Pagination PaginationInfo `json:"pagination"`
		Items      []T            `json:"items"`
	}

	FiltersRequest struct {
		Sort           string `json:"sortBy" query:"sortBy"`
		PageIndex      int    `json:"pageIndex" query:"pageIndex"`
		ResultsPerPage int    `json:"recordsPerPage" query:"recordsPerPage"`
		Search         string `json:"searchText" query:"searchText"`
		IncludeDeleted bool   `json:"includeDeleted" query:"includeDeleted"`
	}

	Filters struct {
		PageIndex             int
		ResultsPerPage        int
		Search                string
		SortColumn            string
		SortDirection         string
		IncludeDeletedRecords bool
	}

	Pagination struct {
		Page           int `json:"page" query:"page"`
		ResultsPerPage int `json:"resultsPerPage" query:"resultsPerPage"`
	}

	// ColumnMeta describes one table column, so a frontend can discover valid
	// sort/filter column names and how to render/parse their values.
	ColumnMeta struct {
		Name     string `json:"name" db:"column_name"`
		DataType string `json:"dataType" db:"data_type"`
	}
)

// TrimmedText trims s and checks it fits a varchar(maxLen) column - at most
// maxLen characters, and non-empty when required - so bad input is a 400
// rather than a database error surfacing as a 500. Returns the trimmed value
// for the caller to store.
func TrimmedText(field string, s string, required bool, maxLen int) (string, error) {
	s = strings.TrimSpace(s)
	if required && s == "" {
		return s, fmt.Errorf("%s is required", field)
	}
	if n := utf8.RuneCountInString(s); n > maxLen {
		return s, fmt.Errorf("%s must be at most %d characters (got %d)", field, maxLen, n)
	}
	return s, nil
}

const (
	DefaultResultsPerPage = 10
	MaxResultsPerPage     = 150
	MaxPageIndex          = 10000
)

// Validate checks the paging values and fills in defaults for omitted ones.
// Negative values would otherwise reach Postgres as a negative LIMIT/OFFSET
// (a 500), and an omitted recordsPerPage would be LIMIT 0 (always empty).
func (f *FiltersRequest) Validate() error {
	if f.ResultsPerPage == 0 {
		f.ResultsPerPage = DefaultResultsPerPage
	}

	if f.ResultsPerPage < 1 || f.ResultsPerPage > MaxResultsPerPage {
		return fmt.Errorf("recordsPerPage(%d) must be between 1 and %d", f.ResultsPerPage, MaxResultsPerPage)
	}

	if f.PageIndex < 0 || f.PageIndex >= MaxPageIndex {
		return fmt.Errorf("pageIndex(%d) must be between 0 and %d", f.PageIndex, MaxPageIndex-1)
	}

	return nil
}
