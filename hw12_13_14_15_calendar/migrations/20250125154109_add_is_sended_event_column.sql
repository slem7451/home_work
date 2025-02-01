-- +goose Up
alter table events add column is_sended boolean default false;

-- +goose Down
alter table events drop column is_sended;
