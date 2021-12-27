CREATE DATABASE blog;

CREATE SCHEMA IF NOT EXISTS blog;   

CREATE TABLE blog.authors (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL
);

CREATE TABLE blog.articles (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	posted_at TIMESTAMP NOT NULL DEFAULT NOW(),
	author_id uuid NOT NULL,
	FOREIGN KEY (author_id)
		REFERENCES blog.authors(id)
		ON DELETE CASCADE
);
