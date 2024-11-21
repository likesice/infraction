create table if not exists infraction(
    id integer,
    name text not null,
    split integer not null,
    created_at integer not null,
    version integer not null,
    primary key(id)
);