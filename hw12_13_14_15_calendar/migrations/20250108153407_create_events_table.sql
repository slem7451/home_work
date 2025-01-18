-- +goose Up
create table events (
    id serial primary key,
    title varchar(255) not null,
    event_date date not null,
    date_since date not null,
    descr text null,
    user_id integer not null,
    notify_date date null
);

-- +goose Down
drop table events;
