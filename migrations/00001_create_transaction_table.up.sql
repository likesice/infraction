create table if not exists "transaction"(
    id integer,
    infraction_id integer not null,
    amount real not null,
    created_at integer not null,
    primary key(id),
    foreign key (infraction_id) references infraction
);