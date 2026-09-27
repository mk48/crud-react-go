-- +migrate Up
CREATE TABLE "user" (
  "id" uuid PRIMARY KEY,
  "sub" varchar(100) UNIQUE NOT NULL,
  "email" varchar(100) UNIQUE NOT NULL,
  "name" varchar(100) NULL,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,

  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "created_by" uuid NOT NULL REFERENCES "user"(id),
  "updated_at" TIMESTAMP WITH TIME ZONE NULL,
  "updated_by" uuid NULL REFERENCES "user"(id),
  "deleted_at" TIMESTAMP WITH TIME ZONE NULL,
  "deleted_by" uuid NULL REFERENCES "user"(id)
);
CREATE INDEX "idx_user_email" ON "user" ("email");

-- Self-referencing system user, used as created_by for seed data.
-- Real "user" rows are otherwise only provisioned on first Casdoor sign-in
-- (see middleware.AuthMiddleware).
INSERT INTO "user" (id, sub, email, name, is_admin, created_at, created_by) VALUES
  ('8e4d4b2a-8dc3-4d73-bfee-7d6032ec2212', 'system', 'kumaran.veera@outlook.com', 'Kumaran', TRUE, now(), '8e4d4b2a-8dc3-4d73-bfee-7d6032ec2212');

CREATE TABLE audit_history (
  id                  UUID PRIMARY KEY,
  source_id           UUID NOT NULL,
  data                JSONB NOT NULL,
  created_at timestamp with time zone NOT NULL DEFAULT (now())
);
CREATE INDEX "idx_audit_history_source_id" ON audit_history(source_id);

-- Sample table for testing the CRUD/auth setup end to end. Not part of the
-- real family-tree schema yet - replace once the actual domain is designed.
CREATE TABLE sample_items (
  id            UUID PRIMARY KEY,
  name          varchar(200) NOT NULL,
  description   TEXT,

  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "created_by" uuid NOT NULL REFERENCES "user"(id),
  "updated_at" TIMESTAMP WITH TIME ZONE NULL,
  "updated_by" uuid NULL REFERENCES "user"(id),
  "deleted_at" TIMESTAMP WITH TIME ZONE NULL,
  "deleted_by" uuid NULL REFERENCES "user"(id)
);

-- +migrate Down
DROP TABLE IF EXISTS sample_items CASCADE;
DROP TABLE IF EXISTS "audit_history" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
