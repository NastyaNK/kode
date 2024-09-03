create table users
(
    id       serial primary key,
    name     varchar(20),
    password varchar(32)
);
create table notes
(
    id      serial primary key,
    user_id int REFERENCES users (id),
    content text
);
