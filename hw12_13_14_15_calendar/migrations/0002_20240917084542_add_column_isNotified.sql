-- +goose Up

ALTER TABLE events ADD COLUMN IF NOT EXISTS is_notified boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE events DROP IF EXISTS COLUMN is_notified;