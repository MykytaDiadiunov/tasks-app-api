CREATE TABLE IF NOT EXISTS users (
    id bigserial not null primary key,
    "name" text not null,
    email text not null,
    avatar text default null,
    password text not null,
    email_confirmed boolean default false,
    email_confirmation_token text
)