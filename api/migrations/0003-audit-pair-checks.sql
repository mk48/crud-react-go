-- +migrate Up
-- Every audited table's *_at/*_by columns are always written together (see
-- the services' Update/Delete and util.UpdateByID). Enforce it, so
-- "deleted_at IS NULL" and "deleted_by IS NULL" are interchangeable
-- soft-delete checks and model.AuditRow.ToDto can rely on both being set.
ALTER TABLE "user"
  ADD CONSTRAINT chk_user_updated_pair CHECK ((updated_at IS NULL) = (updated_by IS NULL)),
  ADD CONSTRAINT chk_user_deleted_pair CHECK ((deleted_at IS NULL) = (deleted_by IS NULL));

ALTER TABLE sample_items
  ADD CONSTRAINT chk_sample_items_updated_pair CHECK ((updated_at IS NULL) = (updated_by IS NULL)),
  ADD CONSTRAINT chk_sample_items_deleted_pair CHECK ((deleted_at IS NULL) = (deleted_by IS NULL));

ALTER TABLE sample_child_items
  ADD CONSTRAINT chk_sample_child_items_updated_pair CHECK ((updated_at IS NULL) = (updated_by IS NULL)),
  ADD CONSTRAINT chk_sample_child_items_deleted_pair CHECK ((deleted_at IS NULL) = (deleted_by IS NULL));

-- +migrate Down
ALTER TABLE sample_child_items
  DROP CONSTRAINT IF EXISTS chk_sample_child_items_updated_pair,
  DROP CONSTRAINT IF EXISTS chk_sample_child_items_deleted_pair;

ALTER TABLE sample_items
  DROP CONSTRAINT IF EXISTS chk_sample_items_updated_pair,
  DROP CONSTRAINT IF EXISTS chk_sample_items_deleted_pair;

ALTER TABLE "user"
  DROP CONSTRAINT IF EXISTS chk_user_updated_pair,
  DROP CONSTRAINT IF EXISTS chk_user_deleted_pair;
