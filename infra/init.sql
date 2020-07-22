\c stone_go

CREATE TABLE IF NOT EXISTS account (
    id serial CONSTRAINT pk_id_account PRIMARY KEY,
    name text NOT NULL,
    cpf text NOT NULL,
    secret text NOT NULL,
    balance numeric(9,2) NOT NULL,
    created_at  timestamp with time zone,
    UNIQUE(cpf)
);

CREATE TABLE IF NOT EXISTS transfer (
    id_transfer serial CONSTRAINT pk_id_transfer PRIMARY KEY,
    acc_origin_id integer NOT NULL,
    acc_destination_id integer NOT NULL,
    amount numeric(9,2) NOT NULL,
    created_at  timestamp with time zone
);


