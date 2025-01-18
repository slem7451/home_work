-- +goose Up
create table events (
    id serial primary key,
    title varchar(255) not null,
    event_date timestamp not null,
    date_since timestamp not null,
    descr text null,
    user_id integer not null,
    notify_date timestamp null
);

-- +goose Down
drop table events;
