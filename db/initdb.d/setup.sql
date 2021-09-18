drop database if exists `sns_db`;
create database `sns_db`;

drop user if exists `sns`@`%`;
create user `sns`@`%` identified by 'sns123456789';
grant all on `sns_db`.* to `sns`@`%`;
flush privileges;

use `sns_db`;

drop table if exists `users`;
drop table if exists `sessions`;

create table `users` (
    `id` serial primary key,
    `name` varchar(255) not null,
    `account_id` varchar(255) not null,
    `password` varchar(255) not null
);

create table `sessions` (
    `id` serial primary key,
    `uuid` text not null,
    `user_id` integer not null,
    `created_at` timestamp not null
);
