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

-- reviews tables
create table kinopoisk.film_reviews
(
    id          bigserial primary key,
    rating      integer,
    body        text,
    product_id  bigint references kinopoisk.films(id) on delete cascade,
    user_id     bigint references kinopoisk.users(id) on delete cascade
);

create table kinopoisk.series_reviews
(
    id          bigserial primary key,
    rating      integer,
    body        text,
    product_id  bigint references kinopoisk.series(id) on delete cascade,
    user_id     bigint references kinopoisk.users(id) on delete cascade
);

-- triggers for review table
create or replace function kinopoisk.film_rating() returns trigger as $film_rating$
	begin
		update kinopoisk.films
			set totalvotes = totalvotes + 1,
				sumvotes = sumvotes + new.rating,
				rating = sumvotes / (totalvotes + 1)
			where id = new.product_id;
		return new;
	end;
$film_rating$ LANGUAGE plpgsql;

create trigger film_rating
	after insert on kinopoisk.film_reviews
	for each row execute procedure kinopoisk.film_rating();
	
create or replace function kinopoisk.series_rating() returns trigger as $series_rating$
	begin
		update kinopoisk.series
			set totalvotes = totalvotes + 1,
				sumvotes = sumvotes + new.rating,
				rating = sumvotes / (totalvotes + 1)
			where id = new.product_id;
		return new;
	end;
$series_rating$ LANGUAGE plpgsql;

create trigger series_rating
	after insert on kinopoisk.series_reviews
	for each row execute procedure kinopoisk.series_rating();

create table kinopoisk.persons
(
    id          bigserial primary key,
    name        varchar(80) not null,
    occupation  varchar(80),
    birth_date  varchar(80),
    birth_place varchar(80)
);

create table kinopoisk.film_actor
(
    id          bigserial primary key,
    film_id     bigint references kinopoisk.films(id) on delete cascade,
    person_id   bigint references kinopoisk.persons(id) on delete cascade
);

create table kinopoisk.series_actor
(
    id          bigserial primary key,
    series_id   bigint references kinopoisk.series(id) on delete cascade,
    person_id   bigint references kinopoisk.persons(id) on delete cascade
);

-- default inserts
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

insert into kinopoisk.series
values (default, 'Комедия', 'Бригада', 'Brigada', '/trailer', 0.12, 10.0, 0, 0, 'brigada description',
        '/static/1.jpg', '/static/1.jpg', 'Россия', 2019, 0, 18),
       (default, 'Комедия', 'Бригада2', 'Brigada2', '/trailer2', 0.12, 10.0, 0, 0, 'brigada2 description',
        '/static/2.jpg', '/static/2.jpg', 'Россия', 2019, 0, 18);

insert into kinopoisk.seasons
values (default, 1, 'season1', 1, 'link1', 'desc1', 2010, 'img1'),
       (default, 1, 'season2', 2, 'link2', 'desc2', 2010, 'img2'),
       (default, 1, 'season3', 3, 'link3', 'desc3', 2011, 'img3'),
       (default, 2, 'season21', 1, 'link21', 'desc21', 2011, 'img21');

insert into kinopoisk.episodes
values (default, 1, 'ep11', 1, 'img1'),
       (default, 1, 'ep12', 2, 'img2'),
       (default, 2, 'ep21', 1, 'img3'),
       (default, 3, 'ep31', 1, 'img3'),
       (default, 4, 'ep41', 1, 'img21');

insert into kinopoisk.persons values(default, 'ivan ivanov', 'actor', '2020-01-01', 'Moscow');
insert into kinopoisk.persons values(default, 'alex alexov', 'actor', '2020-01-01', 'Moscow');

insert into kinopoisk.film_actor values(default, 1, 1);
insert into kinopoisk.series_actor values(default, 1, 2);