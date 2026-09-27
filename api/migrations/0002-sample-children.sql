-- +migrate Up
-- Second sample table, referencing sample_items by foreign key - exercises
-- the "pick a related row" (select/combobox) side of the CRUD setup. Like
-- sample_items, not part of the real family-tree schema.
CREATE TABLE sample_child_items (
  id              UUID PRIMARY KEY,
  sample_item_id  UUID NOT NULL REFERENCES sample_items(id),
  name            varchar(200) NOT NULL,

  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "created_by" uuid NOT NULL REFERENCES "user"(id),
  "updated_at" TIMESTAMP WITH TIME ZONE NULL,
  "updated_by" uuid NULL REFERENCES "user"(id),
  "deleted_at" TIMESTAMP WITH TIME ZONE NULL,
  "deleted_by" uuid NULL REFERENCES "user"(id)
);
CREATE INDEX "idx_sample_child_items_sample_item_id" ON sample_child_items(sample_item_id);

-- +migrate Down
DROP TABLE IF EXISTS sample_child_items CASCADE;
