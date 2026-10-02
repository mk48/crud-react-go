-- +migrate Up

-- Trigram indexes let the list endpoints' free-text search
-- (col ILIKE '%term%', see util.List) use an index instead of scanning.
-- pg_trgm is a trusted extension, so the database owner can create it.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Conventions for every audited table:
--  * *_at/*_by pairs are always set together (CHECK), so "deleted_at IS NULL"
--    alone identifies a live row.
--  * Rows are soft-deleted (deleted_at/deleted_by) and read-only afterwards
--    (see util.UpdateByID).
--  * Default list order is created_at DESC with id as tie-breaker, served by
--    a partial index over live rows.

CREATE TABLE "user" (
  id          uuid PRIMARY KEY,
  sub         varchar(100) NOT NULL UNIQUE,
  -- Stored lower-case (see middleware.AuthMiddleware), so the plain UNIQUE
  -- constraint is effectively case-insensitive.
  email       varchar(254) NOT NULL UNIQUE CHECK (email = lower(email)),
  name        varchar(100) NULL CHECK (name IS NULL OR length(btrim(name)) > 0),
  is_admin    boolean NOT NULL DEFAULT FALSE,
  -- Service accounts act for an API client that has no signed-in user (a
  -- batch job, the API's own internal tasks - see util.ClientRegistry).
  -- They can never sign in interactively (see middleware.AuthMiddleware).
  is_service  boolean NOT NULL DEFAULT FALSE,

  created_at  timestamptz NOT NULL DEFAULT now(),
  created_by  uuid NOT NULL REFERENCES "user"(id),
  updated_at  timestamptz NULL,
  updated_by  uuid NULL REFERENCES "user"(id),
  deleted_at  timestamptz NULL,
  deleted_by  uuid NULL REFERENCES "user"(id),

  CONSTRAINT chk_user_updated_pair CHECK ((updated_at IS NULL) = (updated_by IS NULL)),
  CONSTRAINT chk_user_deleted_pair CHECK ((deleted_at IS NULL) = (deleted_by IS NULL))
);
CREATE INDEX idx_user_live_created_at ON "user" (created_at, id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_name_trgm ON "user" USING gin (name gin_trgm_ops);
CREATE INDEX idx_user_email_trgm ON "user" USING gin (email gin_trgm_ops);

-- Self-referencing system user, used as created_by for seed data.
-- Real "user" rows are otherwise only provisioned on first Casdoor sign-in
-- (see middleware.AuthMiddleware).
INSERT INTO "user" (id, sub, email, name, is_admin, created_at, created_by) VALUES
  ('8e4d4b2a-8dc3-4d73-bfee-7d6032ec2212', 'system', 'kumaran.veera@outlook.com', 'Kumaran', TRUE, now(), '8e4d4b2a-8dc3-4d73-bfee-7d6032ec2212');

-- Service account for work the API does on its own (scheduled/startup
-- tasks - see util.SystemContext), so it's never attributed to a person.
-- Its sub isn't a Casdoor id, so no sign-in can ever resolve to it.
INSERT INTO "user" (id, sub, email, name, is_admin, is_service, created_at, created_by) VALUES
  ('00000000-0000-0000-0000-000000000001', 'service:system', 'system@kfamily.internal', 'System', FALSE, TRUE, now(), '00000000-0000-0000-0000-000000000001');

-- One row per user action that writes data - a form save, an import, a
-- button that updates several tables. Every audit_history row points at the
-- operation that caused it, so a record's history can say *why* it changed
-- (e.g. "created by import of family.csv", or "updated as a side effect of
-- approving sample X") and not just who changed it. See util.RunOperation.
CREATE TABLE operation (
  id            uuid PRIMARY KEY,
  -- Dotted "<resource>.<action>" key, e.g. 'sample.create'. The web app
  -- translates it for display (operation.kind.* in translation.json).
  kind          varchar(100) NOT NULL CHECK (length(btrim(kind)) > 0),
  -- Deferred so a user's own sign-up (middleware.CreateUser) can record its
  -- operation before the "user" row it is performed by exists.
  performed_by  uuid NOT NULL REFERENCES "user"(id) DEFERRABLE INITIALLY DEFERRED,
  -- The record the user acted on directly, if any - audit_history rows of
  -- other records under the same operation are its side effects. NULL for
  -- operations with no single target (e.g. an import).
  target_table  varchar(100) NULL,
  target_id     uuid NULL,
  -- Free-form context, e.g. an import's file name and row count.
  metadata      jsonb NOT NULL DEFAULT '{}',
  -- The app the operation came from: 'web', 'mobile', 'admin',
  -- 'batch:<job>', 'system'. Verified - taken from the access token's
  -- Casdoor application (see util.ClientRegistry), never from a header.
  client        varchar(100) NOT NULL CHECK (length(btrim(client)) > 0),
  -- What the caller reports about itself (app version, platform, IP, user
  -- agent). Useful hints, not proof.
  client_info   jsonb NOT NULL DEFAULT '{}',
  -- W3C trace id (32 lower-case hex) of the request/task that ran it - the
  -- link to its OpenTelemetry trace and logs. NULL only if it ran outside
  -- any trace.
  trace_id      char(32) NULL CHECK (trace_id ~ '^[0-9a-f]{32}$'),
  created_at    timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT chk_operation_target_pair CHECK ((target_table IS NULL) = (target_id IS NULL))
);
-- Default list order (operation.Service.List), newest first.
CREATE INDEX idx_operation_created_at ON operation (created_at, id);
-- "Everything user X did".
CREATE INDEX idx_operation_performed_by ON operation (performed_by, created_at DESC);
-- "Everything done from the mobile app".
CREATE INDEX idx_operation_client ON operation (client, created_at DESC);
-- From a trace (or a log line's trace_id) back to what it changed.
CREATE INDEX idx_operation_trace_id ON operation (trace_id);
CREATE INDEX idx_operation_kind_trgm ON operation USING gin (kind gin_trgm_ops);

-- One row per create/update/delete of any audited table (see
-- util.RecordAudit), holding the full row as it was right after the change.
CREATE TABLE audit_history (
  id           uuid PRIMARY KEY,
  table_name   varchar(100) NOT NULL,
  source_id    uuid NOT NULL,
  action       varchar(10) NOT NULL CHECK (action IN ('create', 'update', 'delete')),
  changed_by   uuid NOT NULL REFERENCES "user"(id),
  -- The operation that caused the change (see util.RunOperation).
  operation_id uuid NOT NULL REFERENCES operation(id),
  data         jsonb NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now()
);
-- A record's history, newest first (audithistory.Service.List).
CREATE INDEX idx_audit_history_source ON audit_history (source_id, created_at DESC);
-- "Everything user X changed" / "all deletes on table Y".
CREATE INDEX idx_audit_history_changed_by ON audit_history (changed_by, created_at DESC);
CREATE INDEX idx_audit_history_table_action ON audit_history (table_name, action, created_at DESC);
-- An operation's changes (operation.Service.GetOne), in the order made.
CREATE INDEX idx_audit_history_operation ON audit_history (operation_id, created_at);

-- Sample table for testing the CRUD/auth setup end to end. Not part of the
-- real family-tree schema yet - replace once the actual domain is designed.
CREATE TABLE sample_items (
  id           uuid PRIMARY KEY,
  name         varchar(200) NOT NULL CHECK (length(btrim(name)) > 0),
  description  text NULL CHECK (description IS NULL OR length(btrim(description)) > 0),

  created_at   timestamptz NOT NULL DEFAULT now(),
  created_by   uuid NOT NULL REFERENCES "user"(id),
  updated_at   timestamptz NULL,
  updated_by   uuid NULL REFERENCES "user"(id),
  deleted_at   timestamptz NULL,
  deleted_by   uuid NULL REFERENCES "user"(id),

  CONSTRAINT chk_sample_items_updated_pair CHECK ((updated_at IS NULL) = (updated_by IS NULL)),
  CONSTRAINT chk_sample_items_deleted_pair CHECK ((deleted_at IS NULL) = (deleted_by IS NULL))
);
CREATE INDEX idx_sample_items_live_created_at ON sample_items (created_at, id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sample_items_name_trgm ON sample_items USING gin (name gin_trgm_ops);
CREATE INDEX idx_sample_items_description_trgm ON sample_items USING gin (description gin_trgm_ops);

-- Second sample table, referencing sample_items by foreign key - exercises
-- the "pick a related row" (select/combobox) side of the CRUD setup. Like
-- sample_items, not part of the real family-tree schema.
CREATE TABLE sample_child_items (
  id              uuid PRIMARY KEY,
  sample_item_id  uuid NOT NULL REFERENCES sample_items(id),
  name            varchar(200) NOT NULL CHECK (length(btrim(name)) > 0),

  created_at      timestamptz NOT NULL DEFAULT now(),
  created_by      uuid NOT NULL REFERENCES "user"(id),
  updated_at      timestamptz NULL,
  updated_by      uuid NULL REFERENCES "user"(id),
  deleted_at      timestamptz NULL,
  deleted_by      uuid NULL REFERENCES "user"(id),

  CONSTRAINT chk_sample_child_items_updated_pair CHECK ((updated_at IS NULL) = (updated_by IS NULL)),
  CONSTRAINT chk_sample_child_items_deleted_pair CHECK ((deleted_at IS NULL) = (deleted_by IS NULL))
);
CREATE INDEX idx_sample_child_items_sample_item_id ON sample_child_items (sample_item_id);
CREATE INDEX idx_sample_child_items_live_created_at ON sample_child_items (created_at, id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sample_child_items_name_trgm ON sample_child_items USING gin (name gin_trgm_ops);

-- +migrate Down
-- No CASCADE: dropping in dependency order fails loudly if anything else
-- was built on these tables, instead of silently dropping it too.
DROP TABLE IF EXISTS sample_child_items;
DROP TABLE IF EXISTS sample_items;
DROP TABLE IF EXISTS audit_history;
DROP TABLE IF EXISTS operation;
DROP TABLE IF EXISTS "user";
