create table if not exists users
(
    id serial not null constraint users_pk primary key,
    wallet_name      varchar(100)               not null,
    wallet_address   varchar(130)               not null,
    "role"           varchar(25) default 'user' not null,
    created          timestamp default now()    not null
);

create table if not exists invitations
(
    referrer_id   int4 references users(id)  on update cascade on delete restrict not null,
    referral_id   int4 references users(id)  on update cascade on delete restrict not null
);

create table if not exists links
(
    id      serial not null constraint links_pk primary key,
    user_id int4 references users(id)  on update cascade on delete restrict not null,
    code    varchar(10) not null
);

create table if not exists boxes
(
    id        serial not null constraint boxes_pk primary key,
    user_id   int4 references users(id)  on update cascade on delete restrict not null,
    available int4 default 0 not null,
    opened    int4 default 0 not null
);

create table if not exists stakes
(
    id        serial not null constraint stakes_pk primary key,
    user_id   int4 references users(id)  on update cascade on delete restrict not null,
    amount    float8                  not null,
    status    boolean   default true  not null,
    tx_hash   varchar(150)            not null,
    created   timestamp default now() not null
);

create table if not exists rewards_pool
(
    id        serial not null constraint rewards_pool_pk primary key,
    available float8 default 0 not null,
    sent      float8 default 0 not null
);

create table if not exists rewards_history
(
    id        serial not null constraint rewards_history_pk primary key,
    user_id   int4 references users(id)  on update cascade on delete restrict not null,
    stake     float8                  not null,
    reward    float8                  not null,
    created   timestamp default now() not null
);

create table if not exists reward_types
(
    id     serial not null constraint reward_types_pk primary key,
    "name" varchar(25) not null
);

INSERT INTO reward_types(id, name) VALUES (1, 'payment');
INSERT INTO reward_types(id, name) VALUES (2, 'box');

create table if not exists rewards
(
    id        serial not null constraint rewards_pk primary key,
    user_id   int4 references users(id)  on update cascade on delete restrict not null,
    status    varchar(25)                  not null,
    type_id   int4 references reward_types(id)  on update cascade on delete restrict not null,
    amount    float8       not null,
    tx_hash   varchar(150) default '',
    created   timestamp default now() not null
);

ALTER TABLE users ADD CONSTRAINT user_wallet_address UNIQUE (wallet_address);
ALTER TABLE links ADD CONSTRAINT links_code UNIQUE (code);
ALTER TABLE stakes ADD CONSTRAINT stakes_tx_hash UNIQUE (tx_hash);
ALTER TABLE reward_types ADD CONSTRAINT reward_types_name UNIQUE ("name");
ALTER TABLE rewards ADD CONSTRAINT rewards_tx_hash UNIQUE (tx_hash);
