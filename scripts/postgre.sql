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


create type filmtype as enum ('films','series');

create table kinopoisk.films
(
    id          serial primary key,
    type        filmtype,
    maingenre   varchar(80), --русский вариант
    russianname varchar(80),
    englishname varchar(80),
    seasons     integer,
    trailerlink varchar(80),
    rating      numeric(4, 2),
    imdbrating  numeric(4, 2),
    description varchar(80), --русское
    image       varchar(80),
    country      varchar(80), --русское
    year        integer,
    agelimit    varchar(80)
);

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

-- example

insert into kinopoisk.films
values (default, 'series', 'Комедия', 'Бригада', 'Brigada', 5, '/trailer', 0.12, 10.0, 'brigada description',
        '/static/1.jpg', 'Россия', 2019, 12),
       (default, 'films', 'Приключения', 'Приключение Мухтара', 'Mukhtar', 1, '/trailer2', 1.23, 9.99,
        'Описание Мухтара', '/static/2.jpg', 'Казахстан', 2017, 6),
       (default, 'films', 'Ужасы', 'Лекарство от здоровья', 'Cure', 1, '/trailer3', 3.45, 4.56,
        'Описание лекарства', '/static/3.jpg', 'США', 2017, 6),
       (default, 'series', 'Приключения', 'Флэш', 'The Flash', 6, '/trailer4', 7.0, 6.9,
        'Описание скорости', '/static/4.jpg', 'США', 2015, 6),
       (default, 'films', 'Приключения', 'Мстители', 'The Avengers', 1, '/trailer5', 9.9, 9.89,
        'Описание мстителей', '/static/5.jpg', 'США', 2015, 6),
       (default, 'series', 'Приключения', 'Время Приключений', 'Adventure Time', 10, '/trailer6', 10.0, 10.0,
        'Джей зе дог энд Фин зе хьюман', '/static/6.jpg', 'США', 2015, 0);

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

select *
from kinopoisk.films f1
         join kinopoisk.films_genres fg1 on (f1.id = fg1.film_id)
         join kinopoisk.genres g1 on (fg1.genre_id = g1.id)
where g1.name = 'Приключения';