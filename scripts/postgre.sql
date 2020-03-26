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

ALTER SEQUENCE kinopoisk.episodes_id_seq RESTART WITH 1; --сериалы ресетить


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

create table kinopoisk.reviews
(
    id       bigserial primary key,
    rating   integer,
    body     text,
    filmId   bigint references kinopoisk.films(id) on delete cascade,
    userId   bigint references kinopoisk.users(id) on delete cascade
);

insert into kinopoisk.films
values (default, 'Комедия', 'Бригада1', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/1.jpg', '/static/1.jpg', 'Россия', 2010, 12),
       (default, 'Приключения', 'Бригада2', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/2.jpg', '/static/2.jpg', 'Россия', 2015, 12),
       (default, 'Приключения', 'Бригада3', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/3.jpg', '/static/3.jpg', 'Россия', 2017, 12),
       (default, 'Ужасы', 'Бригада4', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/4.jpg', '/static/4.jpg', 'Россия', 2019, 12),
       (default, 'Комедия', 'Бригада5', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
        '/static/5.jpg', '/static/5.jpg', 'Россия', 2017, 12),
       (default, 'Ужасы', 'Бригада6', 'Brigada', '/trailer', 0.0, 0.0, 0, 0, 'brigada description',
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

create table kinopoisk.series
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

insert into kinopoisk.series
values (default, 'Комедия', 'Бригада', 'Brigada', '/trailer', 0.12, 10.0, 0, 0, 'brigada description',
        '/static/1.jpg', '/static/1.jpg', 'Россия', 2019, 0, 18),
       (default, 'Комедия', 'Бригада2', 'Brigada2', '/trailer2', 0.12, 10.0, 0, 0, 'brigada2 description',
        '/static/2.jpg', '/static/2.jpg', 'Россия', 2019, 0, 18);

create table kinopoisk.seasons
(
    id          serial primary key,
    seriesid    integer references kinopoisk.series (id),
    name        varchar(80),
    number      integer,
    trailerlink varchar(80),
    description varchar(80),
    year        integer,
    image       varchar(80)
);

insert into kinopoisk.seasons
values (default, 1, 'season1', 1, 'link1', 'desc1', 2010, 'img1'),
       (default, 1, 'season2', 2, 'link2', 'desc2', 2010, 'img2'),
       (default, 1, 'season3', 3, 'link3', 'desc3', 2011, 'img3'),
       (default, 2, 'season21', 1, 'link21', 'desc21', 2011, 'img21');

create table kinopoisk.episodes
(
    id       serial primary key,
    seasonid integer references kinopoisk.seasons (id),
    name     varchar(80),
    number   integer,
    image    varchar(80)
);

insert into kinopoisk.episodes
values (default, 1, 'ep11', 1, 'img1'),
       (default, 1, 'ep12', 2, 'img2'),
       (default, 2, 'ep21', 1, 'img3'),
       (default, 3, 'ep31', 1, 'img3'),
       (default, 4, 'ep41', 1, 'img21');

--example
select *
from kinopoisk.films f1
         join kinopoisk.films_genres fg1 on (f1.id = fg1.film_id)
         join kinopoisk.genres g1 on (fg1.genre_id = g1.id)
where g1.name = 'Приключения';