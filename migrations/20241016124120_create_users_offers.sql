-- +goose Up
-- +goose StatementBegin
create table users (
    id bigint primary key,
    username varchar(32),
    status varchar(20),
    interests varchar[]
);

create table offers (
    user_id bigint unique,
    "text" text,
    interest varchar,
    foreign key (user_id) references users(id) on delete cascade
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists offers;
drop table if exists users;
-- +goose StatementEnd
