create table groups
(
    id integer not null
        primary key autoincrement,
    name varchar(255) not null
        unique
);

create table users
(
    id integer not null
        primary key autoincrement
        unique,
    username varchar(255) not null
        unique,
    password varchar(255),
    name varchar(255) not null,
    locale varchar(255) not null,
    lock_password bool not null,
    blocked bool not null,
    home text not null,
    read int default 0 not null,
    "create" integer default 0 not null,
    modify int default 0 not null,
    "delete" int default 0 not null,
    share int default 0 not null,
    is_admin integer default 0 not null
);

create table user_auth_providers
(
    id varchar(255) not null
        primary key
        unique,
    name varchar(255) not null,
    user_id integer not null
        references users
            on delete set null
);

create table user_groups
(
    user_id varchar(255) not null
        references users
            on delete cascade,
    group_id integer not null
        references groups
            on delete cascade,
    primary key (user_id, group_id)
);

create table volumes
(
    id integer not null
        primary key autoincrement,
    label varchar(255) not null,
    path text not null,
    description text
);

create table group_volumes
(
    group_id integer not null
        references groups
            on delete cascade,
    volume_id integer not null
        references volumes
            on delete cascade,
    read int default 0 not null,
    "create" integer default 0 not null,
    modify int default 0 not null,
    "delete" int default 0 not null,
    share int default 0 not null,
    primary key (group_id, volume_id)
);

create table user_volumes
(
    user_id varchar(255) not null
        references users
            on delete cascade,
    volume_id integer not null
        references volumes
            on delete cascade,
    read int default 0 not null,
    "create" integer default 0 not null,
    modify int default 0 not null,
    "delete" int default 0 not null,
    share int default 0 not null,
    primary key (user_id, volume_id)
);

create index volume_label
    on volumes (label);

create index volume_path
    on volumes (path);

