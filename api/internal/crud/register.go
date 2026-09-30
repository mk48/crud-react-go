package crud

import (
	"context"
	"fmt"
	"log/slog"

	"kfamily/internal/middleware"
	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

// Deps are the shared services Register needs.
type Deps struct {
	DB  *sqlx.DB
	Log *slog.Logger
	MW  *middleware.Middleware
}

// Register mounts res's generic endpoints under api+path, all behind
// AuthMiddleware:
//
//	GET    {path}          List     - page/sort/search (PaginationFilter)
//	GET    {path}/query    Query    - advanced query-builder filter
//	GET    {path}/meta     Meta     - column names and types
//	GET    {path}/:id      Get      - one row (?includeDeleted=true)
//	POST   {path}          Create   - admin only
//	PUT    {path}/:id      Update   - admin only
//	DELETE {path}/:id      Delete   - admin only (soft delete)
//
// minus any in res.Omit. Every write is audited and transactional.
//
// # Custom endpoints
//
// Register returns the route group (already behind AuthMiddleware) and
// the Service, so an entity adds its own endpoints next to the generic
// ones with ordinary Echo handlers in its own package:
//
//	g, svc := crud.Register(ctx, deps, api, "/v1/samples", resource)
//	h := &Handler{svc: svc}
//	g.GET("/summary", h.Summary)                                // any signed-in user
//	g.POST("/:id/archive", h.Archive, deps.MW.AdminMiddleware)   // admins only
//
// Static paths like "/summary" take precedence over "/:id", so they can't
// be shadowed by the generic routes. A custom handler typically:
//
//   - reads the caller with util.CurrentUser(c) and path ids with
//     util.ParseUUIDParam(c, "id");
//   - reuses svc (svc.GetOne, svc.List, ...) or queries svc.DB() directly;
//   - writes through util.WithTx + util.Insert/util.UpdateByID, so the
//     change is audited and atomic;
//   - answers with util.HttpData / util.HttpError like the generic ones;
//   - carries swag annotations (see e.g. user/handler.go's Me) and then
//     `go generate ./cmd/api` to update Swagger.
//
// To replace a generic endpoint rather than add one, list it in
// res.Omit and register your own handler on the same method and path.
// To only tweak one, prefer the hooks (CreateParams, UpdateParams,
// BeforeDelete). User-facing API docs for the generic endpoints live in
// each entity's swagger.go.
func Register[TRow, TDto, TCreate, TUpdate any](
	ctx context.Context,
	deps Deps,
	api *echo.Group,
	path string,
	res Resource[TRow, TDto, TCreate, TUpdate],
) (*echo.Group, *Service[TRow, TDto, TCreate, TUpdate]) {
	if res.Omit&Create == 0 && res.CreateParams == nil {
		panic(fmt.Sprintf("crud: %s needs CreateParams (or Omit: crud.Create)", res.Table))
	}
	if res.Omit&Update == 0 && res.UpdateParams == nil {
		panic(fmt.Sprintf("crud: %s needs UpdateParams (or Omit: crud.Update)", res.Table))
	}

	svc := NewService(deps.DB, res)
	h := &Handler[TRow, TDto, TCreate, TUpdate]{svc: svc}

	g := api.Group(path)
	g.Use(deps.MW.AuthMiddleware)

	// "actionAt"/"actionBy" are the UI's composite audit column (most recent
	// of created/updated) with no single matching DB column - sort by the
	// updated column as a representative approximation.
	paging := deps.MW.PaginationFilter(util.BuildSortableColumns(ctx, deps.Log, deps.DB, res.Table, map[string]string{
		"actionAt": "updated_at",
		"actionBy": "updated_by",
	}))

	has := func(r Route) bool { return res.Omit&r == 0 }
	if has(List) {
		g.GET("", h.List, paging)
	}
	if has(Query) {
		g.GET("/query", h.Query, paging)
	}
	if has(Meta) {
		g.GET("/meta", h.Metadata)
	}
	if has(Get) {
		g.GET("/:id", h.GetOne)
	}
	if has(Create) {
		g.POST("", h.Create, deps.MW.AdminMiddleware)
	}
	if has(Update) {
		g.PUT("/:id", h.Update, deps.MW.AdminMiddleware)
	}
	if has(Delete) {
		g.DELETE("/:id", h.Delete, deps.MW.AdminMiddleware)
	}

	return g, svc
}
