-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS employee (
    id bigint generated always as IDENTITY primary key not null,
    name text not null ,
    create_at timestamptz default now(),
    update_at timestamptz default now()
);

CREATE TABLE IF NOT EXISTS role (
    id bigint generated always as IDENTITY primary key not null,
    name text not null unique,
    create_at timestamptz default now(),
    update_at timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE employee;

DROP TABLE role;
-- +goose StatementEnd
