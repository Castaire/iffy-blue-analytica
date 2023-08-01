CREATE TABLE IF NOT EXISTS albums (
	album_id text NOT NULL PRIMARY KEY,
	album_name text NOT NULL,
	album_uri text NOT NULL,
	album_genres text
);