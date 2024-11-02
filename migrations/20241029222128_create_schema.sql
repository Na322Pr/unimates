-- +goose Up
-- +goose StatementBegin
create table "users" (
  "id" bigint primary key,
  "username" varchar(32) not null,
  "role" varchar(20) default 'user',
  "status" varchar(20) not null,
  "created_at" timestamptz default current_timestamp,
  "modified_at" timestamptz default current_timestamp
);

create table "interests" (
  "id" int primary key,
  "name" varchar(150) unique
);

create table "offers" (
  "id" bigserial primary key,
  "user_id" bigint not null references users(id) on delete cascade,
  "text" text,
  "interest_id" int references interests(id) on delete cascade,
  "notify" boolean default false,
  "inactive_at" timestamptz,
  "created_at" timestamptz default current_timestamp,
  "modified_at" timestamptz default current_timestamp
);

create table "offer_acceptances" (
  "user_id" bigint not null references users(id) on delete cascade,
  "offer_id" bigint not null references offers(id) on delete cascade,
  primary key (user_id, offer_id)
);

create table "user_interests" (
  "user_id" bigint not null references users(id) on delete cascade,
  "interest_id" int not null references interests(id) on delete cascade,
  primary key (user_id, interest_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "user_interests" cascade;
drop table if exists "offer_acceptances" cascade;
drop table if exists "offers" cascade;
drop table if exists "interests" cascade;
drop table if exists "users" cascade;
-- +goose StatementEnd
