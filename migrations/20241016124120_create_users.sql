-- +goose Up
-- +goose StatementBegin
create table users (
    id bigint primary key,
    username varchar(32),
    status varchar(20),
    interests varchar[]
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
