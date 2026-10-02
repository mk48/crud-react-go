# Multiple apps: sign-in and request logging

How each app (web, mobile, admin, batch job) signs in to the API, how the API
knows which app a request comes from, and what gets recorded for each request.
For the full tracing setup, see [tracing.md](tracing.md).

## 1. One Casdoor application per app

Every app that calls the API is its own Casdoor application with its own
`client_id`. The API keeps an allow-list of these applications:

| App | Client name | Registered in | Casdoor grant | Acts as |
|---|---|---|---|---|
| Web app | `web` | `CASDOOR_CLIENT_ID` | authorization code (exchanged by the API) | the signed-in user |
| Mobile app | `mobile` | `API_CLIENTS` | authorization code + PKCE | the signed-in user |
| Admin app | `admin` | `API_CLIENTS` | authorization code + PKCE | the signed-in user |
| Batch job | `batch:<job>` | `API_CLIENTS` with `serviceUserId` | client credentials | its service account |
| API's own tasks | `system` | built in | none | the `system` service account |

```bash
API_CLIENTS=[{"name":"mobile","clientId":"<casdoor client id>"},
             {"name":"batch:nightly-sync","clientId":"<casdoor client id>","serviceUserId":"<user id>"}]
```

Code: [util/client.go](../api/internal/util/client.go) (registry),
[middleware/auth.go](../api/internal/middleware/auth.go) (token checks).

## 2. Signing in

All apps get an access token from Casdoor, then send it on every request as
`Authorization: Bearer <token>`. How they get the token differs:

**Web app.** The browser is redirected to Casdoor's sign-in page and comes back
to `/callback` with a code. The web app sends that code to
`POST /api/v1/auth/signin`, and the API exchanges it for a token using the web
application's client secret. The secret never leaves the server.

**Mobile and admin apps.** They talk to Casdoor directly, using the
authorization code flow with PKCE (no secret in the app). They don't use
`/api/v1/auth/signin`: that endpoint only exchanges codes for the web app's
Casdoor application.

**Batch jobs.** They use the client-credentials grant with their own client id
and secret. There is no user; the job acts as its service account.

## 3. How the API checks each request

`AuthMiddleware` runs on every `/api/v1/*` request (except sign-in):

1. Verify the token's signature with Casdoor's certificate.
2. Check the token type is `access-token` and the issuer is our Casdoor.
3. Find the app: look up the token's `azp` claim (or `aud` if there's no
   `azp`) in the allow-list. Not found means **401**.
4. Store the app name and what the caller says about itself in the request
   context (see section 4).
5. Find the user:
   - **Batch job:** load its service account. If it's deleted, the request
     gets **403**.
   - **Signed-in user:** look them up by Casdoor `sub`, or create them on first
     sign-in. Service accounts can never be signed into, by `sub` or by email
     (**403**).

The app is taken **only from the signed token**. Headers such as
`X-Client-Platform` are recorded, but they can't change which app a request
is attributed to.

## 4. What apps send

| Header | Required | Purpose |
|---|---|---|
| `Authorization: Bearer <token>` | yes | user and app identity |
| `traceparent` | recommended | continues the app's trace into the API |
| `X-Client-Version` | recommended | app version, e.g. `1.4.2` |
| `X-Client-Platform` | recommended | `web`, `ios`, `android`, ... |

The web app sends all four automatically
([api-client.ts](../web/src/lib/api-client.ts)). Every response includes
`X-Trace-Id`, which apps should show in error messages.

## 5. What gets recorded

Every request gets a **trace id**, either the caller's (from `traceparent`) or
a new one. It connects everything below.

| Where | What | Kept |
|---|---|---|
| Response header `X-Trace-Id` | the trace id | n/a |
| Log line (stdout) | method, route, status, latency, IP, user agent, `trace_id`, `span_id` | log retention |
| Trace (OpenTelemetry) | request span plus one span per SQL query; tagged with `kfamily.client`, `kfamily.operation.kind`, `kfamily.operation.id` | tracing backend retention |
| `operation` row (one per write action) | `kind`, `performed_by`, `client`, `client_info` (version, platform, IP, user agent), `trace_id`, target record | permanent |
| `audit_history` rows (one per changed row) | full row snapshot, `changed_by`, `operation_id` | permanent |

Requests that only read data produce a log line and a trace, but no
`operation` row.

Example: a user edits a sample from the mobile app.

```
operation       kind=sample.update  client=mobile  performed_by=<alice>
                client_info={"platform":"ios","version":"2.0.0","ip":"…","userAgent":"…"}
                trace_id=4bf92f3577b34da6a3ce929d0e0e4736
audit_history   table=sample_items  action=update  operation_id=<that operation>
log line        PUT /api/v1/samples/:id 200  trace_id=4bf92f3577b34da6a3ce929d0e0e4736
```

## 6. Finding a request later

Starting from the trace id (e.g. the "Reference" a user saw in an error):

```sql
SELECT o.kind, o.client, o.performed_by, o.created_at, a.table_name, a.source_id, a.action
FROM operation o
JOIN audit_history a ON a.operation_id = o.id
WHERE o.trace_id = '4bf92f3577b34da6a3ce929d0e0e4736';
```

Then search the logs for `trace_id=<id>` and open the trace in the tracing UI.
In the web app, the **Operations** page lists every operation with its source
app; its detail page shows the trace id.

## 7. CORS

Only browser apps on another domain need CORS. Add their origin to
`CORS_ALLOWED_ORIGINS` (explicit list, no `*`):

| App | CORS entry |
|---|---|
| Web app (served by the API) | none, same origin |
| Admin app on its own domain | `https://admin.example.com` |
| Mobile app, batch job | none, not a browser |

## 8. Adding a new app (checklist)

1. Create its Casdoor application (pick the grant from section 1).
2. Batch job only: create its service account (`"user".is_service = TRUE`).
   For now this is done with SQL; see [tracing.md](tracing.md#add-a-batch-job).
3. Add it to `API_CLIENTS`. The API checks the list at startup.
4. Browser app on another domain only: add its origin to
   `CORS_ALLOWED_ORIGINS`.
5. In the app, send the headers from section 4 and show `X-Trace-Id` on errors.

No API code changes are needed. From then on, its writes are recorded with
its client name.
