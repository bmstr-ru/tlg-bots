create table if not exists lot
(
    id           text        not null primary key,
    status       text        not null,
    name         text        not null,
    description  text        not null,
    create_date  timestamptz not null,
    bid_end_time timestamptz not null,
    image        bytea       null,
    is_stopped   boolean,
    has_appeals  boolean,
    is_annulled  boolean
);