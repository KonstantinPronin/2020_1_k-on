CREATE DATABASE k_on
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'C.UTF8'
    LC_CTYPE = 'C.UTF8'
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

create table kinopoisk.genres
(
--     id        serial primary key,
    name      varchar(80)        not null, --русское
    reference varchar(80) unique not null  --англ
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

create table kinopoisk.films_genres
(
    film_id   int references kinopoisk.films (id) on update cascade on delete cascade,
    genre_ref varchar(80) references kinopoisk.genres (reference) on update cascade on delete cascade,
    constraint films_genres_pkey primary key (film_id, genre_ref)
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
    genre_ref varchar(80) references kinopoisk.genres (reference) on update cascade on delete cascade,
    constraint series_genres_pkey primary key (series_id, genre_ref)
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

create table kinopoisk.playlists
(
    id      bigserial primary key,
    name    varchar(80) not null,
    public  bool default false,
    user_id bigint references kinopoisk.users (id) on delete cascade,
    unique (name, user_id)
);

create table kinopoisk.film_playlist
(
    id          bigserial primary key,
    playlist_id bigint references kinopoisk.playlists (id) on delete cascade,
    film_id     bigint references kinopoisk.films (id) on delete cascade,
    unique (playlist_id, film_id)
);

create table kinopoisk.series_playlist
(
    id          bigserial primary key,
    playlist_id bigint references kinopoisk.playlists (id) on delete cascade,
    series_id   bigint references kinopoisk.series (id) on delete cascade,
    unique (playlist_id, series_id)
);

create table kinopoisk.subscriptions
(
    id          bigserial primary key,
    playlist_id bigint references kinopoisk.playlists (id) on delete cascade,
    user_id     bigint references kinopoisk.users (id) on delete cascade,
    unique (playlist_id, user_id)
);

alter table kinopoisk.films
    add column textsearchable_index_col tsvector;

update kinopoisk.films
SET textsearchable_index_col =
        to_tsvector('russian', coalesce(russianname, '') || ' ' || coalesce(description, ''));

create index films_textsearchable_idx on kinopoisk.films using gin (textsearchable_index_col);

create or replace function kinopoisk.film_searchable_text() returns trigger as
$film_searchable_text$
begin
    update kinopoisk.films
    set textsearchable_index_col =
            to_tsvector('russian', coalesce(new.russianname, '') || ' ' || coalesce(new.description, ''))
    where id = new.id;
    return new;
end;
$film_searchable_text$ LANGUAGE plpgsql;

create trigger film_searchable_text
    after insert
    on kinopoisk.films
    for each row
execute procedure kinopoisk.film_searchable_text();

alter table kinopoisk.series
    add column textsearchable_index_col tsvector;

update kinopoisk.series
SET textsearchable_index_col =
        to_tsvector('russian', coalesce(russianname, '') || ' ' || coalesce(description, ''));

create index series_textsearchable_idx on kinopoisk.series using gin (textsearchable_index_col);

create or replace function kinopoisk.series_searchable_text() returns trigger as
$series_searchable_text$
begin
    update kinopoisk.series
    set textsearchable_index_col =
            to_tsvector('russian', coalesce(new.russianname, '') || ' ' || coalesce(new.description, ''))
    where id = new.id;
    return new;
end;
$series_searchable_text$ LANGUAGE plpgsql;

create trigger series_searchable_text
    after insert
    on kinopoisk.series
    for each row
execute procedure kinopoisk.series_searchable_text();

alter table kinopoisk.persons
    add column textsearchable_index_col tsvector;

update kinopoisk.persons
SET textsearchable_index_col =
        to_tsvector('russian', coalesce("name", ''));

create index persons_textsearchable_idx on kinopoisk.persons using gin (textsearchable_index_col);

create or replace function kinopoisk.persons_searchable_text() returns trigger as
$persons_searchable_text$
begin
    update kinopoisk.persons
    set textsearchable_index_col =
            to_tsvector('russian', coalesce(new.name, ''))
    where id = new.id;
    return new;
end;
$persons_searchable_text$ LANGUAGE plpgsql;

create trigger persons_searchable_text
    after insert
    on kinopoisk.persons
    for each row
execute procedure kinopoisk.persons_searchable_text();

select f1.russianname, count(fp2.film_id)
from kinopoisk.film_playlist fp1
         join kinopoisk.film_playlist fp2 on fp1.playlist_id = fp2.playlist_id
         join kinopoisk.films f1 on fp2.film_id = f1.id
where fp1.film_id = 1
group by f1.russianname
order by count(fp2.film_id) desc;

select f2.*
from kinopoisk.films f2
         join (select f1.russianname, count(fp2.film_id)
               from kinopoisk.film_playlist fp1
                        join kinopoisk.film_playlist fp2 on fp1.playlist_id = fp2.playlist_id
                        join kinopoisk.films f1 on fp2.film_id = f1.id
               where fp1.film_id = 1
               group by f1.russianname) as sub on f2.russianname = sub.russianname
order by sub.count desc
offset 1;
