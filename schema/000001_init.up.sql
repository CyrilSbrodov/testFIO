CREATE TABLE users (
    id bigserial not null primary key,
    firstname varchar not null,
    lastname varchar not null,
    patronymic varchar,
    age integer,
    gender varchar,
    nation varchar
    );