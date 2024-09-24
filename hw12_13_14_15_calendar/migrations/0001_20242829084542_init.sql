-- +goose Up
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; - must be run on a database with proper rights

CREATE table events
(
    id                  uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title               varchar(64) not null,
    duration            tsrange not null,
    description         text,
    user_id             uuid      not null,
    notify_before_event bigint
);

-- +goose Down
drop table events;