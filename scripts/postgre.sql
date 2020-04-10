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
        'Отважному путешественнику Питеру Квиллу попадает в руки таинственный артефакт, принадлежащий могущественному ' ||
        'и безжалостному злодею Ронану, строящему коварные планы по захвату Вселенной. Питер оказывается в центре ' ||
        'межгалактической охоты, где жертва — он сам.Единственный способ спасти свою жизнь — объединиться с ' ||
        'четверкой нелюдимых изгоев: воинственным енотом по кличке Ракета, человекоподобным деревом Грутом, ' ||
        'смертельно опасной Гаморой и одержимым жаждой мести Драксом, также известным как Разрушитель. '
           , '/static/img/A1.jpg', '/static/img/A01.jpg', 'Россия', 2010, 13),
       (default, 'horrors', 'Очень страшное кино', 'Scary Movie', 'nIW4y4w502M', 0.0, 0.0, 0, 0,
        'Дрю Деккер разговаривает по телефону с незнакомцем, будто бы попавшим не туда. Очень скоро она понимает, ' ||
        'что её собеседник — маньяк, который хочет её убить. Дрю выбегает из дома, по пятам преследуемая маньяком. ' ||
        'Девушке почти удаётся убежать, но её случайно сбивает машина отца…',
        '/static/img/A2.jpg', '/static/img/A20.jpg', 'Россия', 2015, 13),
       (default, 'war', 'Т-34', 'T-34', 'DrtRCIQVHnA', 0.0, 0.0, 0, 0,
        'В ноябре 1941 года в деревне Нефёдовка вчерашний курсант Николай Ивушкин был вынужден принять командование ' ||
        'единственным уцелевшим танком и вступить в неравное противостояние с немецкой танковой ротой, приближавшейся ' ||
        'к Москве. Силами противника командовал опытный гауптман Клаус Ягер, но это не спасло его подразделение от ' ||
        'поражения, хоть самому ему и удалось ранить нашего курсанта.Летом 1944 года уже пытавшийся совершить ' ||
        'несколько побегов Ивушкин попадет в концлагерь в Тюрингии. Там же оказывается и его «старый знакомый» ' ||
        'Ягер, который решает использовать навыки Николая для тренировки немецких кадетов. Он распоряжается собрать ' ||
        'тому команду и выдать только что прибывшую с фронта новую модель танка Т-34, не зная, что идею побега Ивушкин не оставил.',
        '/static/img/A3.jpg', '/static/img/A30.jpg', 'Россия', 2017, 13),
       (default, 'historical', 'Принц Персии', 'Prince of Persia', 'uOCYsMJHmKk', 0.0, 0.0, 0, 0,
        'Главный герой этой экранизации культовой компьютерной игры, юный принц Дастан всегда побеждал врагов в бою, ' ||
        'но потерял королевство из-за козней коварного царедворца. Теперь Дастану предстоит похитить из рук злодеев ' ||
        'могущественный магический артефакт, способный повернуть время вспять и сделать своего владельца властелином мира.',
        '/static/img/A4.jpg', '/static/img/A40.jpg', 'Россия', 2019, 14),
       (default, 'animated', 'Твое имя', 'Kimi no Na wa', 'tT7b5wR0IOM', 0.0, 0.0, 0, 0,
        'История о парне из Токио и девушке из провинции, которые обнаруживают, что между ними существует странная ' ||
        'и необъяснимая связь. Во сне они меняются телами и проживают жизни друг друга. Но однажды эта способность ' ||
        'исчезает так же внезапно, как появилась. Таки решает во что бы то ни стало отыскать Мицуху.',
        '/static/img/A5.jpg', '/static/img/A50.jpg', 'Россия', 2017, 11),
       (default, 'detectives', 'Убийство в восточном экспрессе', 'Murder on the Orient Express', 'pTK0hUqzolU', 0.0,
        0.0, 0, 0, 'Путешествие на одном из самых роскошных поездов Европы неожиданно превращается в одну из ' ||
                   'самых стильных и захватывающих загадок в истории. Фильм рассказывает историю тринадцати пассажиров ' ||
                   'поезда, каждый из которых находится под подозрением. И только сыщик должен как можно быстрее ' ||
                   'разгадать головоломку, прежде чем преступник нанесет новый удар.',
        '/static/img/A6.jpg', '/static/img/60.jpg', 'Россия', 2017, 11),
       (default, 'biographical', 'Господин Никто', 'Mr. Nobody', 'aAWGDN2S-sE', 0.0, 0.0, 0, 0,
        'Проснувшийся немощным стариком Немо Никто оказывается последним смертным в гротескном будущем. ' ||
        'Все люди уже давно бессмертны и с удовольствием наблюдают за телешоу, где главная звезда — дряхлый и ' ||
        'безумный старик Немо, доживающий свои последние дни. Накануне конца к нему приходит журналист, и Немо ' ||
        'рассказывает ему свою историю перескакивая из одной жизни в другую, параллельную, ' ||
        'несколько раз за рассказ успев умереть.',
        '/static/img/A7.jpg', '/static/img/A70.jpg', 'Россия', 2017, 18),
       (default, 'documentary', 'Он вам не Димон', 'He is not Dimon to You', 'qrwlk7_GF9g', 0.0, 0.0, 0, 0,
        'В фильме рассказывается о предполагаемом недвижимом имуществе председателя Правительства Российской ' ||
        'Федерации Дмитрия Медведева.', '/static/img/A8.jpg', '/static/img/A80.jpg', 'Россия', 2017, 18),
       (default, 'criminal', 'Джокер', 'Joker', 'jGfiPs9zuhE', 0.0, 0.0, 0, 0,
        'Готэм, начало 1980-х годов. Комик Артур Флек живет с больной матерью, которая с детства учит ' ||
        'его «ходить с улыбкой». Пытаясь нести в мир хорошее и дарить людям радость, Артур сталкивается с ' ||
        'человеческой жестокостью и постепенно приходит к выводу, что этот мир получит от него ' ||
        'не добрую улыбку, а ухмылку злодея Джокера.',
        '/static/img/A9.jpg', '/static/img/A90.jpg', 'Россия', 2017, 18),
       (default, 'action', 'Бригада', 'Brigada', 'l3F5Tu1AZUU', 0.0, 0.0, 0, 0,
        'Это история четырех друзей детства, обычных московских парней, Саши Белого, Космоса, Пчелы и Фила, ' ||
        'выросших в одном дворе. Друзья решили немного подзаработать, но незапланированное убийство вмиг перемешало все ' ||
        'задуманное, поставив на кон их жизни.Ставка слишком высока, но отступать некуда. Теперь парни прокладывают ' ||
        'себе дорогу в криминальном мире и волею судеб превращаются в одну из самых сплоченных и влиятельных группировок…',
        '/static/img/A10.jpg', '/static/img/A100.jpg', 'Россия', 2017, 18),
       (default, 'drama', 'Титаник', 'Titanic', 'ZQ6klONCq4s', 0.0, 0.0, 0, 0,
        'В первом и последнем плавании шикарного «Титаника» встречаются двое. Пассажир нижней палубы Джек выиграл билет ' ||
        'в карты, а богатая наследница Роза отправляется в Америку, чтобы выйти замуж по расчёту. Чувства молодых людей ' ||
        'только успевают расцвести, и даже не классовые различия создадут испытания влюблённым, а айсберг, ' ||
        'вставший на пути считавшегося непотопляемым лайнера.',
        '/static/img/A11.jpg', '/static/img/A110.jpg', 'Россия', 2017, 12),
       (default, 'melodrama', 'Три метра над уровнем неба', 'Three Steps Above The Sky', 'elLW9pxj4kM', 0.0, 0.0, 0, 0,
        'История двух молодых людей, которые принадлежат к разным мирам. Баби ― богатая девушка, которая отображает доброту' ||
        ' и невинность. Аче — мятежный мальчик, импульсивный, бессознательный, склонный к риску и опасности. ' ||
        'Это маловероятно, практически невозможно, но их встреча неизбежна, и в этом неистовом путешествии ' ||
        'между ними возникает первая большая любовь.', '/static/img/A12.jpg', '/static/img/120.jpg', 'Россия', 2017,
        18),
       (default, 'comedy', 'Третий лишний', 'TED', 'tJIcP7YGhw0', 0.0, 0.0, 0, 0,
        'Джон влюблен в красавицу Лори. Он работает в прокате автомобилей и имеет большие планы на будущее. ' ||
        'Но в их отношения вмешивается третий, давний друг Джона — Тед. Он отрывается сутки напролет, предпочитает ' ||
        'случайные связи и не желает терять друга. Но никто на самом деле не знает на что он способен, ' ||
        'ведь Тед — большой плюшевый медведь.',
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
        'Действие фэнтезийного мультсериала «Время приключений» разворачивается в постапокалиптическом мире под ' ||
        'названием Земля Ууу, пережившем ядерную войну. С тех пор прошло почти тысячу лет, все выжившие люди ' ||
        'подверглись мутациям, и в мир пришла магия. В центре сюжета тринадцатилетний подросток по имени Финн ' ||
        'и его верный приятель пес по кличке Джейк, который ко всему прочему является для мальчика еще и наставником. ' ||
        'Они живут вместе в доме на дереве, и больше всего на свете любят разные приключения, которые либо настигают ' ||
        'парочку сами, либо им подкидают жители Земли Ууу. Финн и Джейк самые настоящие герои, ведь они всегда ' ||
        'готовы прийти на помощь всем кто попал в беду. Бороться со злом друзьям по большей части помогают ' ||
        'волшебные способности Джейка, тело которого способно растягиваться до любых размеров и форм',
        '/static/img/1.jpg', '/static/img/01.jpg', 'Россия', 2010, 2011, 6),
       (default, 'horrors', 'Американская история ужасов', 'American Horror Story', '6Vw-xG8nLyE', 0.0, 0.0, 0, 0,
        'Американская история ужасов – это сериал, где каждый сезон рассказывает новую историю, пропитанную атмосферой ' ||
        'страха, трепета и безысходности. Многосерийное шоу изобилует паранормальными явлениями, мистикой, а многие ' ||
        'сюжеты основаны на реальных событиях. Среди рассказанных историй присутствуют: дом с призраками, ужасы ' ||
        'закрытой психиатрической лечебницы, школа ведьм из Нового Орлеана, цирковое шоу с уродами, загадка ' ||
        'таинственного отеля «Кортез», злоключения межрасовой супружеской пары, мистика президентских выборов ' ||
        'в США 2016 года и даже вымирание человечества из-за апокалипсиса.', '/static/img/2.jpg',
        '/static/img/20.jpg', 'Россия', 2015, 2016, 16),
       (default, 'war', 'Братья по оружию', 'Band of Brothers', 'panok5dLHM4', 0.0, 0.0, 0, 0,
        'Американская история ужасов – это сериал, где каждый сезон рассказывает новую историю, пропитанную ' ||
        'атмосферой страха, трепета и безысходности. Многосерийное шоу изобилует паранормальными явлениями, мистикой, ' ||
        'а многие сюжеты основаны на реальных событиях. Среди рассказанных историй присутствуют: дом с призраками, ' ||
        'ужасы закрытой психиатрической лечебницы, школа ведьм из Нового Орлеана, цирковое шоу с уродами, загадка ' ||
        'таинственного отеля «Кортез», злоключения межрасовой супружеской пары, мистика президентских выборов в ' ||
        'США 2016 года и даже вымирание человечества из-за апокалипсиса.',
        '/static/img/3.jpg', '/static/img/30.jpg', 'Россия', 2017, 2018, 16),
       (default, 'historical', 'Чернобыль', 'Chernobyl', 'qtY2sel76qo', 0.0, 0.0, 0, 0,
        'В сериале рассказывается о боевом пути роты E («Easy») 2-го батальона 506-го парашютно-десантного полка ' ||
        '101-й воздушно-десантной дивизии США от тренировочного лагеря в Таккоа, штат Джорджия, через высадку в ' ||
        'Нормандии, операцию «Маркет Гарден» и Бастонское сражение до конца войны.',
        '/static/img/4.jpg', '/static/img/40.jpg', 'Россия', 2019, 2020, 16),
       (default, 'animated', 'Наруто', 'Naruto', 'xTOJ5_RKdl8', 0.0, 0.0, 0, 0,
        'Это история, в которой рассказывается про мальчика-ниндзя. Он мечтает стать Хокаге: главой своей деревни. ' ||
        'Но Хокаге – это самый мудрый и сильный ниндзя деревни, поэтому парень попытается преодолеть кучу испытаний, ' ||
        'победить множество противников, заслужить уважение, подрасти морально и физически.',
        '/static/img/5.jpg', '/static/img/50.jpg', 'Россия', 2017, 2018, 8),
       (default, 'detectives', 'Шерлок', 'Sherlock', 'eMM7sX4-6gc', 0.0, 0.0, 0, 0,
        'События разворачиваются в наши дни. Он прошел Афганистан, остался инвалидом. По возвращении в родные ' ||
        'края встречается с загадочным, но своеобразным гениальным человеком. Тот в поиске соседа по квартире. ' ||
        'Лондон, 2010 год. Происходят необъяснимые убийства. Скотланд-Ярд без понятия, за что хвататься. ' ||
        'Существует лишь один человек, который в силах разрешить проблемы и найти ответы на сложные вопросы.',
        '/static/img/6.jpg', '/static/img/60.jpg', 'Россия', 2017, 2018, 12),
       (default, 'biographical', 'Высоцкий.Четыре часа настоящей жизни', 'Vysotsky', 'wTbqwQbLmOA', 0.0, 0.0, 0, 0,
        'Сериал рассказывает об одном из наиболее драматичных периодов жизни Владимира Высоцкого. ' ||
        'Мы застаем героя на пике известности. В конце 70-х Высоцкий — самый знаменитый человек в СССР, ' ||
        'кумир миллионов. Но его силы на исходе. Он измотан, он больше не ощущает свой дар, он не может писать стихи…',
        '/static/img/7.jpg', '/static/img/70.jpg', 'Россия', 2017, 2018, 16),
       (default, 'documentary', 'Как устроена наша планета', 'How our Earth is build', 'SgApoHS6eJE', 0.0, 0.0, 0, 0,
        'Это история о создании всего в этом мире. Программа исследует, как Вселенная возникла из ничего, ' ||
        'и как она выросла с точки значительно меньше, чем атомные частицы, до огромного космоса.',
        '/static/img/8.jpg', '/static/img/80.jpg', 'Россия', 2017, 2018, 16),
       (default, 'criminal', 'Острые козырьки', 'Peaky Blinders', '0InieEzg5kY', 0.0, 0.0, 0, 0,
        'Британский сериал о криминальном мире Бирмингема 20-х годов прошлого века, в котором многолюдная ' ||
        'семья Шелби стала одной из самых жестоких и влиятельных гангстерских банд послевоенного времени. ' ||
        'Фирменным знаком группировки, промышлявшей грабежами и азартными играми, стали зашитые в козырьки лезвия.',
        '/static/img/9.jpg', '/static/img/90.jpg', 'Россия', 2017, 2018, 16),
       (default, 'action', 'Флеш', 'The Flash', 'rNfU1myyZYo', 0.0, 0.0, 0, 0,
        'Когда Барри Аллен был маленьким, больше всего на свете ему хотелось быть супергероем – тем, ' ||
        'кто превосходит лимиты человеческого организма и использует во благо данную ему силу. Когда ' ||
        'Барри было одиннадцать лет, он на собственном опыте узнал, что люди с необычными способностями ' ||
        'действительно существуют: его мать была убита одним из таких людей. Повзрослев и став судмедэкспертом, ' ||
        'Барри не отбросил мыслей о сверхлюдях и продолжал искать доказательства их существования, что не лучшим ' ||
        'образом сказывалось на его репутации и общении с коллегами. Впрочем, однажды его усилия были вознаграждены...',
        '/static/img/10.jpg', '/static/img/100.jpg', 'Россия', 2017, 2018, 16),
       (default, 'drama', 'Сверхъестественное', 'Supernatural', 'la_XCx06Ric', 0.0, 0.0, 0, 0,
        'Сериал рассказывает о приключениях братьев Сэма и Дина Винчестеров, которые путешествуют по ' ||
        'Соединённым Штатам на чёрном автомобиле Chevrolet Impala 1967 года, расследуют паранормальные ' ||
        'явления, многие из которых основаны на американских городских легендах и фольклоре, и сражаются с ' ||
        'порождениями зла, такими как демоны и призраки.',
        '/static/img/11.jpg', '/static/img/110.jpg', 'Россия', 2017, 2018, 18),
       (default, 'melodrama', 'Однажды в сказке', 'Once Upon A Time', 'JV4v-Lu3NUI', 0.0, 0.0, 0, 0,
        'Сюжет фэнтези разворачивается в двух мирах - современном и сказочном. Жизнь 28-летней Эммы Свон ' ||
        'меняется, когда ее 10-летний сын Генри, от которого она отказалась много лет назад, находит Эмму и ' ||
        'объявляет, что она является дочерью Прекрасного Принца и Белоснежки. Само собой разумеется, у мальчишки ' ||
        'нет никаких сомнений, что параллельно нашему существует альтернативный сказочный мир – город Сторибрук, в ' ||
        'котором в итоге оказывается Эмма. Постепенно героиня привязывается к необычному мальчику и странному ' ||
        'городу, жители которого «забыли», кем они были в прошлом. А все из-за проклятия Злой Королевы (по ' ||
        'совместительству приемной матери Генри), с помощью которого колдунья остановила время в сказочной стране. ' ||
        'Однако стоит протянуть руку - и сказка оживет. Эпическая битва за будущее двух миров начинается, но, чтобы ' ||
        'одержать победу, Эмме придется принять свою судьбу...',
        '/static/img/120.jpg', '/static/img/12.jpg', 'Россия', 2017, 2018, 18),
       (default, 'comedy', 'Теория большого взрыва', 'The Big Bang Theory', '48p90POPC5I', 0.0, 0.0, 0, 0,
        'Два блестящих физика Леонард и Шелдон - великие умы, которые понимают, как устроена вселенная. ' ||
        'Но их гениальность ничуть не помогает им общаться с людьми, особенно с женщинами. Всё начинает меняться, ' ||
        'когда напротив них поселяется красавица Пенни.Стоит также отметить пару странных друзей этих физиков: ' ||
        'Воловиц который любит употреблять фразы на разных языках, включая русский, а Раджеш Кутрапали теряет дар ' ||
        'речи при виде женщин.',
        '/static/img/13.jpg', '/static/img/130.jpg', 'Россия', 2017, 2018, 18);

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

update kinopoisk.persons
set image = '/static/img/person1.jpg'
where id = 1;
update kinopoisk.persons
set image = '/static/img/person2.jpg'
where id = 2;
update kinopoisk.persons
set image = '/static/img/person3.jpg'
where id = 3;