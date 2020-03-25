create table kinopoisk.users
(
    id       bigserial primary key,
    username varchar(80) not null,
    password varchar(80) not null,
    email    varchar(80),
    image    varchar(80)
);

grant all privileges on schema kinopoisk to usr;
grant all privileges on all tables in schema kinopoisk to usr;
grant all privileges on all sequences in schema kinopoisk to usr;

create table kinopoisk.films
(
    id              serial primary key,
    maingenre       varchar(80), --русский вариант
    russianname     varchar(80),
    englishname     varchar(80),
    trailerlink     varchar(80),
    rating          numeric(4, 2),
    imdbrating      numeric(4, 2),
    totalvotes      integer,
    sumvotes        integer,
    description     varchar(80), --русское
    image           varchar(80),
    backgroundimage varchar(80),
    country         varchar(80), --русское
    year            integer,
    agelimit        varchar(80)
);

insert into kinopoisk.films
values (default, 'Комедия', 'Бригада', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/1.jpg', '/static/1.jpg', 'Россия', 2010, 12),
       (default, 'Приключения', 'Бригада', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/2.jpg', '/static/2.jpg', 'Россия', 2015, 12),
       (default, 'Приключения', 'Бригада', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/3.jpg', '/static/3.jpg', 'Россия', 2017, 12),
       (default, 'Ужасы', 'Бригада', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/4.jpg', '/static/4.jpg', 'Россия', 2019, 12),
       (default, 'Комедия', 'Бригада', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/5.jpg', '/static/5.jpg', 'Россия', 2017, 12),
       (default, 'Ужасы', 'Бригада', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/6.jpg', '/static/6.jpg', 'Россия', 2017, 12);

create table kinopoisk.genres
(
    id        serial primary key,
    name      varchar(80) not null, --русское
    reference varchar(80) not null  --англ
);

create table kinopoisk.films_genres
(
    film_id  int references kinopoisk.films (id) on update cascade on delete cascade,
    genre_id int references kinopoisk.genres (id) on update cascade on delete cascade,
    constraint films_genres_pkey primary key (film_id, genre_id)
);

insert into kinopoisk.genres
values (default, 'Приключения', 'Adventures'),
       (default, 'Ужасы', 'Horros'),
       (default, 'Комедия', 'Comedy');

insert into kinopoisk.films_genres
values (1, 3),
       (1, 1),
       (2, 1),
       (3, 2),
       (3, 1),
       (4, 1),
       (4, 3),
       (5, 1),
       (6, 1),
       (6, 3);

create table kinopoisk.serials
(
    id              serial primary key,
    maingenre       varchar(80), --русский вариант
    russianname     varchar(80),
    englishname     varchar(80),
    trailerlink     varchar(80),
    rating          numeric(4, 2),
    imdbrating      numeric(4, 2),
    totalvotes      integer,
    sumvotes        integer,
    description     varchar(80), --русское
    image           varchar(80),
    backgroundimage varchar(80),
    country         varchar(80), --русское
    yearfirst       integer,
    yearlast        integer,
    agelimit        varchar(80)
);

insert into kinopoisk.serials
values (default, 'Комедия', 'Бригада', 'Brigada', '/trailer', 0.12, 10.0, 'brigada description',
        '/static/1.jpg', 'Россия', 2019, 12);


--example
select *
from kinopoisk.films f1
         join kinopoisk.films_genres fg1 on (f1.id = fg1.film_id)
         join kinopoisk.genres g1 on (fg1.genre_id = g1.id)
where g1.name = 'Приключения';