-- -- +goose Up
-- -- +goose StatementBegin
do $$
declare
    max_id integer;
begin
    select coalesce(max(id), 0) into max_id from interests;
	execute format('create sequence interests_id_seq start with %s minvalue %s', max_id + 1, max_id + 1);
	execute 'alter table interests alter column id set default nextval(''interests_id_seq'')';
    execute 'alter sequence interests_id_seq owned by interests.id';
end $$;
-- -- +goose StatementEnd

-- -- +goose Down
-- -- +goose StatementBegin
-- SELECT 'down SQL query';
drop sequence if exists interests_id_seq cascade;
-- -- +goose StatementEnd