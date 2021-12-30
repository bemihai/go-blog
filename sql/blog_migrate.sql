-- CREATE EXTENSION pgcrypto;

CREATE TABLE authors (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL
);

CREATE TABLE articles (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	posted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	author_id uuid NOT NULL,
	FOREIGN KEY (author_id)
		REFERENCES authors(id)
		ON DELETE CASCADE
);

CREATE INDEX ON "authors" ("id");

CREATE INDEX ON "articles" ("id");

CREATE INDEX ON "articles" ("author_id");
