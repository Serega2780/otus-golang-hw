-- +goose Up
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; - must be run on a database with proper wrights

CREATE table events
(
    id                  uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title               varchar(64) not null,
    duration            tsrange not null,
    description         text,
    user_id             uuid      not null,
    notify_before_event bigint
);

INSERT INTO events (title, duration, description, user_id, notify_before_event)
VALUES ('1st event', '[2024-08-28T10:00:00Z, 2024-08-28T11:00:00Z]', 'long description for 1st event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', null),
       ('2nd event', '[2024-08-28T12:00:00Z, 2024-08-28T13:00:00Z]', 'long description for 2nd event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', null),
       ('3rd event', '[2024-08-28T14:00:00Z, 2024-08-28T15:00:00Z]', 'long description for 3rd event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', null);

-- +goose Down
drop table events;