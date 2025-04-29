-- +goose Up
-- +goose StatementBegin
create table public.users
(
    u_id bigserial not null,
    email text not null,
    pass_hash bytea not null,
    primary key (u_id)
);
create table public.apps
(
    a_id bigserial not null,
    name text not null,
    secret_key text not null,
    primary key(a_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table public.users;
drop table public.apps;
-- +goose StatementEnd
