package user

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
// @Summary      Get one user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id              path   string  true   "user id (UUID)"
// @Param        includeDeleted  query  bool    false  "Return it even if soft-deleted"
// @Success      200  {object}  util.HttpResult{result=user.Dto}
// @Failure      400,401,404,500  {object}  util.HttpResult
// @Router       /api/v1/users/{id} [get]
func swaggerGetOne() {}

// swaggerList documents the generic crud List endpoint.
//
// @Summary      List users
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  One page of users, optionally narrowed by a free-text search.
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        searchText      query  string  false  "Free-text search"
// @Param        includeDeleted  query  bool    false  "Include soft-deleted rows"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[user.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/users [get]
func swaggerList() {}

// swaggerMetadata documents the generic crud Metadata endpoint.
//
// @Summary      Users table columns
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Column names and Postgres data types - valid sort/filter columns for the query builder.
// @Success      200  {object}  util.HttpResult{result=[]dto.ColumnMeta}
// @Failure      401,500  {object}  util.HttpResult
// @Router       /api/v1/users/meta [get]
func swaggerMetadata() {}

// swaggerQuery documents the generic crud Query endpoint.
//
// @Summary      Advanced query for users
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Rows matching a parameterized SQL WHERE fragment from the query builder: quoted column names and :p1, :p2... placeholders, validated server-side. Includes soft-deleted rows (filter on deleted_at to exclude them).
// @Param        pageIndex       query  int     false  "Page index, from 0"  default(0)
// @Param        recordsPerPage  query  int     false  "Rows per page (1-150)"  default(10)
// @Param        sortBy          query  string  false  "Sort as <column>:asc|desc, e.g. createdAt:desc"
// @Param        whereCondition                query  string  true   "WHERE fragment with double-quoted column names and :pN placeholders"
// @Param        whereConditionParametersJson  query  string  false  "JSON array of the :p1, :p2... values, in order"
// @Success      200  {object}  util.HttpResult{result=dto.PaginationResponse[user.Dto]}
// @Failure      400,401,500  {object}  util.HttpResult
// @Router       /api/v1/users/query [get]
func swaggerQuery() {}

// swaggerUpdate documents the generic crud Update endpoint.
//
// @Summary      Update a user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft-deleted rows can't be updated (404). Omit isAdmin to leave it unchanged; removing the last admin is refused (409).
// @Accept       json
// @Param        id       path  string  true  "user id (UUID)"
// @Param        request  body  user.UpdateInputDto  true  "New values"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,409,500  {object}  util.HttpResult
// @Router       /api/v1/users/{id} [put]
func swaggerUpdate() {}

// swaggerDelete documents the generic crud Delete endpoint.
//
// @Summary      Delete a user
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Description  Admin only. Soft delete - the row and its audit history are kept. You can't delete yourself or the last admin (409).
// @Param        id  path  string  true  "user id (UUID)"
// @Success      200  {object}  util.HttpResult
// @Failure      400,401,403,404,409,500  {object}  util.HttpResult
// @Router       /api/v1/users/{id} [delete]
func swaggerDelete() {}
