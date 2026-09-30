package samplechild

// API docs for the generic crud endpoints (see crud.Register). swag reads
// route docs from function comments and the handlers are generic, so they
// live here, on documentation-only functions. Regenerate the spec after
// editing: go generate ./cmd/api

import (
	"kfamily/internal/dto"
	"kfamily/internal/util"
)

// swag resolves the dto./util. type names below through this file's imports.
var _ = []any{dto.ColumnMeta{}, util.HttpResult{}}

// swaggerGetOne documents the generic crud GetOne endpoint.
//
// @Summary      Get one sample child item
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id              path   string  true   "sample child item id (UUID)"
// @Param        includeDeleted  query  bool    false  "Return it even if soft-deleted"
// @Success      200  {object}  util.HttpResult{result=samplechild.Dto}
// @Failure      400,401,404,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children/{id} [get]
func swaggerGetOne() {}

// swaggerList documents the generic crud List endpoint.
//
// @Summary      List sample child items
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Description  One page of sample child items, optionally narrowed by a free-text search.
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        searchText      query  string  false  "Free-text search"
// @Param        includeDeleted  query  bool    false  "Include soft-deleted rows"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[samplechild.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children [get]
func swaggerList() {}

// swaggerMetadata documents the generic crud Metadata endpoint.
//
// @Summary      Sample children table columns
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Column names and Postgres data types - valid sort/filter columns for the query builder.
// @Success      200  {object}  util.HttpResult{result=[]dto.ColumnMeta}
// @Failure      401,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children/meta [get]
func swaggerMetadata() {}

// swaggerQuery documents the generic crud Query endpoint.
//
// @Summary      Advanced query for sample child items
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Rows matching a parameterized SQL WHERE fragment from the query builder: quoted column names and :p1, :p2... placeholders, validated server-side. Includes soft-deleted rows (filter on deleted_at to exclude them).
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        whereCondition                query  string  true   "WHERE fragment with double-quoted column names and :pN placeholders"
// @Param        whereConditionParametersJson  query  string  false  "JSON array of the :p1, :p2... values, in order"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[samplechild.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children/query [get]
func swaggerQuery() {}

// swaggerCreate documents the generic crud Create endpoint.
//
// @Summary      Create a sample child item
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only.
// @Accept       json
// @Param        request  body  samplechild.CreateInputDto  true  "New sample child item"
// @Success      201  {object}  util.HttpResult{result=string}  "id of the new sample child item"
// @Failure      400,401,403,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children [post]
func swaggerCreate() {}

// swaggerUpdate documents the generic crud Update endpoint.
//
// @Summary      Update a sample child item
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft-deleted rows can't be updated (404).
// @Accept       json
// @Param        id       path  string  true  "sample child item id (UUID)"
// @Param        request  body  samplechild.UpdateInputDto  true  "New values"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children/{id} [put]
func swaggerUpdate() {}

// swaggerDelete documents the generic crud Delete endpoint.
//
// @Summary      Delete a sample child item
// @Tags         Sample children
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft delete - the row and its audit history are kept.
// @Param        id  path  string  true  "sample child item id (UUID)"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,500  {object}  util.HttpResult
// @Router       /api/v1/sample-children/{id} [delete]
func swaggerDelete() {}
