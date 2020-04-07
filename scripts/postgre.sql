CREATE DATABASE k_on
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'C'
    LC_CTYPE = 'C'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

GRANT ALL ON DATABASE k_on TO k_on;

GRANT ALL ON DATABASE k_on TO postgres;

GRANT TEMPORARY, CONNECT ON DATABASE k_on TO PUBLIC;

CREATE SCHEMA kinopoisk
    AUTHORIZATION postgres;

create table kinopoisk.users
(
    id       bigserial primary key,
    username varchar(80) not null,
    password varchar(80) not null,
    email    varchar(80),
    image    varchar(80)
);

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

create table kinopoisk.episodes
(
    id       serial primary key,
    seasonid integer references kinopoisk.seasons (id),
    name     varchar(80),
    number   integer,
    image    varchar(80)
);

create table kinopoisk.series_genres
(
    series_id int references kinopoisk.series (id) on update cascade on delete cascade,
    genre_id  int references kinopoisk.genres (id) on update cascade on delete cascade,
    constraint series_genres_pkey primary key (series_id, genre_id)
);

-- reviews tables
create table kinopoisk.film_reviews
(
    id         bigserial primary key,
    rating     integer,
    body       text,
    product_id bigint references kinopoisk.films (id) on delete cascade,
    user_id    bigint references kinopoisk.users (id) on delete cascade
);

create table kinopoisk.series_reviews
(
    id         bigserial primary key,
    rating     integer,
    body       text,
    product_id bigint references kinopoisk.series (id) on delete cascade,
    user_id    bigint references kinopoisk.users (id) on delete cascade
);

-- triggers for review table
create or replace function kinopoisk.film_rating() returns trigger as
$film_rating$
begin
    update kinopoisk.films
    set totalvotes = totalvotes + 1,
        sumvotes   = sumvotes + new.rating,
        rating     = (sumvotes + new.rating) / (totalvotes + 1)
    where id = new.product_id;
    return new;
end;
$film_rating$ LANGUAGE plpgsql;

create trigger film_rating
    after insert
    on kinopoisk.film_reviews
    for each row
execute procedure kinopoisk.film_rating();

create or replace function kinopoisk.series_rating() returns trigger as
$series_rating$
begin
    update kinopoisk.series
    set totalvotes = totalvotes + 1,
        sumvotes   = sumvotes + new.rating,
        rating     = (sumvotes + new.rating) / (totalvotes + 1)
    where id = new.product_id;
    return new;
end;
$series_rating$ LANGUAGE plpgsql;

create trigger series_rating
    after insert
    on kinopoisk.series_reviews
    for each row
execute procedure kinopoisk.series_rating();

create table kinopoisk.persons
(
    id          bigserial primary key,
    name        varchar(80) not null,
    occupation  varchar(80),
    birth_date  varchar(80),
    birth_place varchar(80)
);

ALTER TABLE kinopoisk.persons
    ADD COLUMN image character varying(80);

create table kinopoisk.film_actor
(
    id        bigserial primary key,
    film_id   bigint references kinopoisk.films (id) on delete cascade,
    person_id bigint references kinopoisk.persons (id) on delete cascade
);

create table kinopoisk.series_actor
(
    id        bigserial primary key,
    series_id bigint references kinopoisk.series (id) on delete cascade,
    person_id bigint references kinopoisk.persons (id) on delete cascade
);

-- default inserts
insert into kinopoisk.films
values (default, 'adventures', 'Стражи галактики', 'Guardians of the Galaxy', 'nlysXqG-gbQ', 0.0, 0.0, 0, 0,
        'Description', '/static/img/A1.jpg', '/static/img/A01.jpg', 'Россия', 2010, 13),
       (default, 'horrors', 'Очень страшное кино', 'Scary Movie', 'nIW4y4w502M', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A2.jpg', '/static/img/A20.jpg', 'Россия', 2015, 13),
       (default, 'war', 'Т-34', 'T-34', 'DrtRCIQVHnA', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A3.jpg', '/static/img/A30.jpg', 'Россия', 2017, 13),
       (default, 'historical', 'Принц Персии', 'Prince of Persia', 'uOCYsMJHmKk', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A4.jpg', '/static/img/A40.jpg', 'Россия', 2019, 14),
       (default, 'animated', 'Твое имя', 'Kimi no Na wa', 'tT7b5wR0IOM', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A5.jpg', '/static/img/A50.jpg', 'Россия', 2017, 11),
       (default, 'detectives', 'Убийство в восточном экспрессе', 'Murder on the Orient Express', 'pTK0hUqzolU', 0.0,
        0.0, 0, 0, 'Description', '/static/img/A6.jpg', '/static/img/60.jpg', 'Россия', 2017, 11),
       (default, 'biographical', 'Господин Никто', 'Mr. Nobody', 'aAWGDN2S-sE', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A7.jpg', '/static/img/A70.jpg', 'Россия', 2017, 18),
       (default, 'documentary', 'Он вам не Димон', 'He is not Dimon to You', 'qrwlk7_GF9g', 0.0, 0.0, 0, 0,
        'Description', '/static/img/A8.jpg', '/static/img/A80.jpg', 'Россия', 2017, 18),
       (default, 'criminal', 'Джокер', 'Joker', 'jGfiPs9zuhE', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A9.jpg', '/static/img/A90.jpg', 'Россия', 2017, 18),
       (default, 'action', 'Бригада', 'Brigada', 'l3F5Tu1AZUU', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A10.jpg', '/static/img/A100.jpg', 'Россия', 2017, 18),
       (default, 'drama', 'Титаник', 'Titanic', 'ZQ6klONCq4s', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A11.jpg', '/static/img/A110.jpg', 'Россия', 2017, 12),
       (default, 'melodrama', 'Три метра над уровнем неба', 'Three Steps Above The Sky', 'elLW9pxj4kM', 0.0, 0.0, 0, 0,
        'Description', '/static/img/A12.jpg', '/static/img/120.jpg', 'Россия', 2017, 18),
       (default, 'comedy', 'Третий лишний', 'TED', 'tJIcP7YGhw0', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/A13.jpg', '/static/img/A130.jpg', 'Россия', 2017, 18);

insert into kinopoisk.genres
values (default, 'Приключения', 'adventures'),
       (default, 'Ужасы', 'horros'),
       (default, 'Военные', 'war'),
       (default, 'Исторические', 'historical'),
       (default, 'Анимация', 'animated'),
       (default, 'Детективы', 'detectives'),
       (default, 'Биографические', 'biographical'),
       (default, 'Документальные', 'documentary'),
       (default, 'Криминал', 'criminal'),
       (default, 'Боевики', 'action'),
       (default, 'Драмы', 'drama'),
       (default, 'Мелодрамы', 'melodrama'),
       (default, 'Комедии', 'comedy');

insert into kinopoisk.films_genres
values (1, 1),
       (2, 2),
       (3, 3),
       (4, 4),
       (5, 5),
       (6, 6),
       (7, 7),
       (8, 8),
       (9, 9),
       (10, 10),
       (11, 11),
       (12, 12),
       (13, 13);

insert into kinopoisk.series
values (default, 'adventures', 'Время приключений', 'Adventure Time', '594sVuwYTKQ', 0.0, 0.0, 0, 0,
        'Description', '/static/img/1.jpg', '/static/img/01.jpg', 'Россия', 2010, 2011, 6),
       (default, 'horrors', 'Американская история ужасов', 'American Horror Story', '6Vw-xG8nLyE', 0.0, 0.0, 0, 0,
        'Description', '/static/img/2.jpg', '/static/img/20.jpg', 'Россия', 2015, 2016, 16),
       (default, 'war', 'Братья по оружию', 'Band of Brothers', 'panok5dLHM4', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/3.jpg', '/static/img/30.jpg', 'Россия', 2017, 2018, 16),
       (default, 'historical', 'Чернобыль', 'Chernobyl', 'qtY2sel76qo', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/4.jpg', '/static/img/40.jpg', 'Россия', 2019, 2020, 16),
       (default, 'animated', 'Наруто', 'Naruto', 'xTOJ5_RKdl8', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/5.jpg', '/static/img/50.jpg', 'Россия', 2017, 2018, 8),
       (default, 'detectives', 'Шерлок', 'Sherlock', 'eMM7sX4-6gc', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/6.jpg', '/static/img/60.jpg', 'Россия', 2017, 2018, 12),
       (default, 'biographical', 'Высоцкий.Четыре часа настоящей жизни', 'Vysotsky', 'wTbqwQbLmOA', 0.0, 0.0, 0, 0,
        'Description', '/static/img/7.jpg', '/static/img/70.jpg', 'Россия', 2017, 2018, 16),
       (default, 'documentary', 'Как устроена наша планета', 'How our Earth is build', 'SgApoHS6eJE', 0.0, 0.0, 0, 0,
        'Description', '/static/img/8.jpg', '/static/img/80.jpg', 'Россия', 2017, 2018, 16),
       (default, 'criminal', 'Острые козырьки', 'Peaky Blinders', '0InieEzg5kY', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/9.jpg', '/static/img/90.jpg', 'Россия', 2017, 2018, 16),
       (default, 'action', 'Флеш', 'The Flash', 'rNfU1myyZYo', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/10.jpg', '/static/img/100.jpg', 'Россия', 2017, 2018, 16),
       (default, 'drama', 'Сверхъестественное', 'Supernatural', 'la_XCx06Ric', 0.0, 0.0, 0, 0, 'Description',
        '/static/img/11.jpg', '/static/img/110.jpg', 'Россия', 2017, 2018, 18),
       (default, 'melodrama', 'Однажды в сказке', 'Once Upon A Time', 'JV4v-Lu3NUI', 0.0, 0.0, 0, 0,
        'Description', '/static/img/120.jpg', '/static/img/12.jpg', 'Россия', 2017, 2018, 18),
       (default, 'comedy', 'Теория большого взрыва', 'The Big Bang Theory', '48p90POPC5I', 0.0, 0.0, 0, 0,
        'Description', '/static/img/13.jpg', '/static/img/130.jpg', 'Россия', 2017, 2018, 18);

insert into kinopoisk.series_genres
values (1, 1),
       (2, 2),
       (3, 3),
       (4, 4),
       (5, 5),
       (6, 6),
       (7, 7),
       (8, 8),
       (9, 9),
       (10, 10),
       (11, 11),
       (12, 12),
       (13, 13);

insert into kinopoisk.seasons
values (default, 1, 'season1', 1, 'link1', 'desc1', 2010, 'img1'),
       (default, 2, 'season2', 2, 'link2', 'desc2', 2010, 'img2'),
       (default, 3, 'season3', 3, 'link3', 'desc3', 2011, 'img3'),
       (default, 4, 'season21', 1, 'link21', 'desc21', 2011, 'img21');

insert into kinopoisk.episodes
values (default, 1, 'ep11', 1, 'img1'),
       (default, 1, 'ep12', 2, 'img2'),
       (default, 2, 'ep21', 1, 'img3'),
       (default, 3, 'ep31', 1, 'img3'),
       (default, 4, 'ep41', 1, 'img21');

insert into kinopoisk.persons
values (default, 'ivan ivanov', 'actor', '2020-01-01', 'Moscow');
insert into kinopoisk.persons
values (default, 'alex alexov', 'actor', '2020-01-01', 'Moscow');

insert into kinopoisk.film_actor
values (default, 1, 1);
insert into kinopoisk.film_actor
values (default, 2, 1);
insert into kinopoisk.film_actor
values (default, 3, 1);
insert into kinopoisk.film_actor
values (default, 4, 1);
insert into kinopoisk.film_actor
values (default, 5, 1);
insert into kinopoisk.film_actor
values (default, 6, 1);
insert into kinopoisk.film_actor
values (default, 7, 1);
insert into kinopoisk.film_actor
values (default, 8, 1);
insert into kinopoisk.film_actor
values (default, 9, 1);
insert into kinopoisk.film_actor
values (default, 10, 1);
insert into kinopoisk.film_actor
values (default, 11, 1);
insert into kinopoisk.film_actor
values (default, 12, 1);
insert into kinopoisk.film_actor
values (default, 13, 1);

insert into kinopoisk.film_actor
values (default, 1, 2);
insert into kinopoisk.film_actor
values (default, 2, 2);
insert into kinopoisk.film_actor
values (default, 3, 2);
insert into kinopoisk.film_actor
values (default, 4, 2);
insert into kinopoisk.film_actor
values (default, 5, 2);
insert into kinopoisk.film_actor
values (default, 6, 2);
insert into kinopoisk.film_actor
values (default, 7, 2);
insert into kinopoisk.film_actor
values (default, 8, 2);
insert into kinopoisk.film_actor
values (default, 9, 2);
insert into kinopoisk.film_actor
values (default, 10, 2);
insert into kinopoisk.film_actor
values (default, 11, 2);
insert into kinopoisk.film_actor
values (default, 12, 2);
insert into kinopoisk.film_actor
values (default, 13, 2);


insert into kinopoisk.series_actor
values (default, 1, 2);
insert into kinopoisk.series_actor
values (default, 2, 2);
insert into kinopoisk.series_actor
values (default, 3, 2);
insert into kinopoisk.series_actor
values (default, 4, 2);
insert into kinopoisk.series_actor
values (default, 5, 2);
insert into kinopoisk.series_actor
values (default, 6, 2);
insert into kinopoisk.series_actor
values (default, 7, 2);
insert into kinopoisk.series_actor
values (default, 8, 2);
insert into kinopoisk.series_actor
values (default, 9, 2);
insert into kinopoisk.series_actor
values (default, 10, 2);
insert into kinopoisk.series_actor
values (default, 11, 2);
insert into kinopoisk.series_actor
values (default, 12, 2);
insert into kinopoisk.series_actor
values (default, 13, 2);

insert into kinopoisk.series_actor
values (default, 1, 3);
insert into kinopoisk.series_actor
values (default, 2, 3);
insert into kinopoisk.series_actor
values (default, 3, 3);
insert into kinopoisk.series_actor
values (default, 4, 3);
insert into kinopoisk.series_actor
values (default, 5, 3);
insert into kinopoisk.series_actor
values (default, 6, 3);
insert into kinopoisk.series_actor
values (default, 7, 3);
insert into kinopoisk.series_actor
values (default, 8, 3);
insert into kinopoisk.series_actor
values (default, 9, 3);
insert into kinopoisk.series_actor
values (default, 10, 3);
insert into kinopoisk.series_actor
values (default, 11, 3);
insert into kinopoisk.series_actor
values (default, 12, 3);
insert into kinopoisk.series_actor
values (default, 13, 3);

update kinopoisk.persons set image = '/static/img/person1.jpg' where id = 1;
update kinopoisk.persons set image = '/static/img/person2.jpg' where id = 2;
update kinopoisk.persons set image = '/static/img/person3.jpg' where id = 3;