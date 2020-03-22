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

create table kinopoisk.genres
(
    id        serial primary key,
    name      varchar(80) not null,
    reference varchar(80) not null
);

create table kinopoisk.films_genres
(
    film_id  int references kinopoisk.films (id) on update cascade on delete cascade,
    genre_id int references kinopoisk.genres (id) on update cascade on delete cascade,
    constraint films_genres_pkey primary key (film_id, genre_id)
);

-- example
select *
from kinopoisk.films f1
         join kinopoisk.films_genres fg1 on (f1.id = fg1.film_id)
         join kinopoisk.genres g1 on (fg1.genre_id = g1.id)
where g1.name = 'drama';