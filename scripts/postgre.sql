create table kinopoisk.users (
	id bigserial primary key,
	username varchar(80) not null,
	password varchar(80) not null,
	email varchar(80),
	image varchar(80)
);

