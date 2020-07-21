\c stone_go
create table if not exists account (
    id serial not null primary key,
    name text not null,
    cpf text unique not null,
    secret text not null,
    balance numeric(9,2) not null,
    created_at  timestamp with time zone
);

create table if not exists transfer (
    id serial not null primary key,
    acc_origin_id integer not null,
    acc_destination_id integer not null,
    amount numeric(9,2) not null,
    created_at  timestamp with time zone
);
