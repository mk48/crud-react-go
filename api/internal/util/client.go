package util

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Built-in client names. Every other app is registered in API_CLIENTS (see
// NewClientRegistry).
const (
	// The web app this binary serves - CASDOOR_CLIENT_ID's application.
	ClientWeb = "web"
	// The API itself (scheduled/startup tasks) - see SystemContext.
	ClientSystem = "system"
)

// Headers a client sends about itself - recorded in operation.client_info.
const (
	ClientVersionHeader  = "X-Client-Version"  // e.g. "1.4.2"
	ClientPlatformHeader = "X-Client-Platform" // e.g. "web", "ios", "android"
)

// SystemUserID is the "system" service account seeded by
// migrations/0001-init.sql - the performed_by of everything SystemContext
// runs.
var SystemUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

var clientNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9:_-]{0,99}$`)

// Client is one app allowed to call the API - a Casdoor application. It's
// identified by the access token it presents, never by anything the caller
// says about itself.
type Client struct {
	// Recorded as operation.client, e.g. "mobile" or "batch:nightly-sync".
	Name string `json:"name"`
	// The Casdoor application's client id - matched against the token's
	// azp/aud claim.
	ClientID string `json:"clientId"`
	// Set for clients with no signed-in user (batch jobs using the
	// client-credentials grant): the service account ("user" row with
	// is_service) they act as.
	ServiceUserID *uuid.UUID `json:"serviceUserId,omitempty"`
}

// IsService reports whether c acts as a service account rather than for a
// signed-in user.
func (c Client) IsService() bool {
	return c.ServiceUserID != nil
}

// ClientRegistry is the allow-list of apps that may call the API.
type ClientRegistry struct {
	byClientID map[string]Client
}

// NewClientRegistry registers the web app (webClientID, from
// CASDOOR_CLIENT_ID) plus every app in apiClientsJson (API_CLIENTS, a JSON
// array of Client), e.g.
//
//	[{"name":"mobile","clientId":"..."},
//	 {"name":"batch:nightly-sync","clientId":"...","serviceUserId":"..."}]
func NewClientRegistry(webClientID string, apiClientsJson string) (*ClientRegistry, error) {
	clients := []Client{{Name: ClientWeb, ClientID: webClientID}}
	if apiClientsJson != "" {
		var extra []Client
		if err := json.Unmarshal([]byte(apiClientsJson), &extra); err != nil {
			return nil, fmt.Errorf("API_CLIENTS is not a valid JSON array of clients: %w", err)
		}
		clients = append(clients, extra...)
	}

	r := &ClientRegistry{byClientID: map[string]Client{}}
	names := map[string]bool{}
	for _, c := range clients {
		if !clientNamePattern.MatchString(c.Name) {
			return nil, fmt.Errorf("client name %q must be 1-100 of a-z, 0-9, ':', '_', '-'", c.Name)
		}
		if c.Name == ClientSystem {
			return nil, fmt.Errorf("client name %q is reserved for the API's own tasks", c.Name)
		}
		if c.ClientID == "" {
			return nil, fmt.Errorf("client %q has no clientId", c.Name)
		}
		if names[c.Name] {
			return nil, fmt.Errorf("client name %q is registered twice", c.Name)
		}
		if _, ok := r.byClientID[c.ClientID]; ok {
			return nil, fmt.Errorf("clientId of %q is registered twice", c.Name)
		}
		names[c.Name] = true
		r.byClientID[c.ClientID] = c
	}

	return r, nil
}

// Lookup returns the client registered for a Casdoor application's client
// id.
func (r *ClientRegistry) Lookup(clientID string) (Client, bool) {
	c, ok := r.byClientID[clientID]
	return c, ok
}

// CheckServiceUsers fails if a service client's serviceUserId isn't a live
// service account - run at startup, so a typo surfaces before the first
// batch run does.
func (r *ClientRegistry) CheckServiceUsers(ctx context.Context, db *sqlx.DB) error {
	for _, c := range r.byClientID {
		if !c.IsService() {
			continue
		}
		var isService bool
		err := db.GetContext(ctx, &isService, `SELECT is_service FROM "user" WHERE id = $1 AND deleted_at IS NULL`, *c.ServiceUserID)
		if errors.Is(err, sql.ErrNoRows) || (err == nil && !isService) {
			return fmt.Errorf("client %q: serviceUserId %s is not a live service account (\"user\".is_service)", c.Name, c.ServiceUserID)
		} else if err != nil {
			return fmt.Errorf("client %q: unable to check its service account. err: %w", c.Name, err)
		}
	}
	return nil
}

// RequestSource is where a write comes from, recorded on its operation (see
// RunOperation).
type RequestSource struct {
	// The verified client name (see Client.Name, ClientSystem).
	Client string
	// What the caller reports about itself - app version, platform, IP,
	// user agent. Stored as operation.client_info; hints, not proof.
	Info map[string]any
}

type requestSourceCtxKey struct{}

// WithRequestSource returns ctx carrying src - set by AuthMiddleware for
// requests and by SystemContext for the API's own tasks.
func WithRequestSource(ctx context.Context, src RequestSource) context.Context {
	return context.WithValue(ctx, requestSourceCtxKey{}, src)
}

func requestSource(ctx context.Context) (RequestSource, bool) {
	src, ok := ctx.Value(requestSourceCtxKey{}).(RequestSource)
	return src, ok
}

// SystemContext starts work the API does on its own - a scheduled or
// startup task named task - as a new trace, attributed to the "system"
// client. Run its writes with PerformedBy: SystemUserID, and end the span
// when the task is done:
//
//	ctx, span := util.SystemContext(ctx, "cleanup-expired-invites")
//	defer span.End()
func SystemContext(ctx context.Context, task string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("kfamily").Start(ctx, task, trace.WithNewRoot())
	return WithRequestSource(ctx, RequestSource{
		Client: ClientSystem,
		Info:   map[string]any{"task": task},
	}), span
}
