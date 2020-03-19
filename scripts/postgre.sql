create table kinopoisk.users
(
    id       bigserial primary key,
    username varchar(80) not null,
    password varchar(80) not null,
    email    varchar(80),
    image    varchar(80)
);

create type filmtype as enum ('film','serial');

create table kinopoisk.films
(
    id          serial primary key,
    type        filmtype,
    maingenre   varchar(80),
    russianname varchar(80),
    englishname varchar(80),
    seasons     integer,
    trailerlink varchar(80),
    rating      numeric(4, 2),
    imdbrating  numeric(4, 2),
    description varchar(80),
    image       varchar(80),
    county      varchar(80),
    year        integer,
    agelimit    varchar(80)
);

