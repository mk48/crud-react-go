package dto

import (
	"fmt"
	"time"
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

	// IdNameBilingual is IdName for a reference table with separate Danish/
	// English names (e.g. place_types) instead of a single Name.
	IdNameBilingual struct {
		Id     string `json:"id"`
		NameDa string `json:"nameDa"`
		NameEn string `json:"nameEn"`
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

func (f *FiltersRequest) Validate() error {
	if f.ResultsPerPage >= 150 {
		return fmt.Errorf("Results/page(%d) is exceed the limit 150", f.ResultsPerPage)
	}

	if f.PageIndex >= 10000 {
		return fmt.Errorf("PageIndex(%d) must be less than 10000", f.PageIndex)
	}

	return nil
}
