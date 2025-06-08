CREATE TABLE role (
    id bigint generated always as IDENTITY primary key not null,
    name text not null unique,
    create_at timestamptz default now(),
    update_at timestamptz default now()
);