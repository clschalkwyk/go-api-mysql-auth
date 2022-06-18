create database if not exists authdb;
use authdb;
create table users (
                       id  bigint(11) unsigned auto_increment primary key,
                       email varchar(30) not null,
                       password varchar(120) not null,
                       salt varchar(120) not null,
                       created_at timestamp default current_timestamp,
                       last_login datetime
)
    charset = utf8;
create index users_email_idx on users (email);
