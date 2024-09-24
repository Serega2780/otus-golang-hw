-- +goose Up
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; - must be run on a database with proper rights

CREATE table events
(
    id                  uuid PRIMARY KEY     DEFAULT uuid_generate_v4(),
    title               varchar(64) not null,
    duration            tsrange     not null,
    description         text,
    user_id             uuid        not null,
    notify_before_event bigint,
    is_notified         boolean     NOT NULL DEFAULT false
);

INSERT INTO events (title, duration, description, user_id, notify_before_event)
VALUES ('1st event', tsrange((now() + interval '15m')::timestamp, (now() + interval '45m')::timestamp), 'long description for 1st event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', 1200000000000),
       ('2nd event', tsrange((now() + interval '1h')::timestamp, (now() + interval '2h')::timestamp), 'long description for 2nd event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', 1200000000000),
       ('3rd event', tsrange((now() + interval '3h')::timestamp, (now() + interval '4h')::timestamp), 'long description for 3rd event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', 1200000000000),
       ('4th event', tsrange((now() + interval '5d')::timestamp, (now() + interval '5d 1h')::timestamp), 'long description for 4th event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', 1200000000000),
       ('5th event', tsrange((now() + interval '20d')::timestamp, (now() + interval '20d 1h')::timestamp), 'long description for 5th event',
        '76dffd8b-4699-47a7-abb4-0a99fa4e8b3b', 1200000000000);

-- +goose Down
drop table events;