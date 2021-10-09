-- drop database if exists "sns_db";
-- create database "sns_db";

-- drop user if exists "sns_user"@"%";
-- create user "sns_user"@"%" identified by 'sns123456789';
-- grant all on "sns_db".* to "sns_user"@"%";
-- flush privileges;

-- use "sns_db";

drop table if exists "users";
drop table if exists "sessions";
drop table if exists "cigarettes";
drop table if exists "follows";

create table "users" (
    "id" serial primary key,
    "name" varchar(255) not null,
    "account_id" varchar(255) unique not null,
    "password" varchar(255) not null
);

create table "sessions" (
    "id" serial primary key,
    "uuid" text not null unique,
    "user_id" bigint not null,
    "created_at" timestamp not null
);

create table "cigarettes" (
    "id" serial primary key,
    "smoked_count" int not null,
    "user_id" bigint not null,
    "created_at" timestamp not null
);

create table "follows" (
    "id" serial primary key,
    "src_account_id" varchar(255) not null,
    "dst_account_id" varchar(255) not null
);
