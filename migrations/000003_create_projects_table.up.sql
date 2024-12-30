CREATE TABLE IF NOT EXISTS projects (
    id bigserial not null primary key,
    title text not null,
    description text,
    creator_id bigserial not null,
    constraint fk_creator
        foreign key (creator_id)
        references users(id)
        on delete cascade
)
