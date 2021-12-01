create database blog;

create schema if not exists blog;   

create table blog.authors (
	id serial,
	name varchar(36) not null,
	email varchar(36),
	primary key (id)
);

create table blog.articles (
	id serial,
	title varchar(255) not null,
	body varchar(255) not null,
	posted_at timestamp not null default now(),
	author_id int not null,
	primary key (id),
	foreign key (author_id)
		references authors(id)
		on delete cascade
);
