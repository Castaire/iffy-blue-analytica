CREATE TABLE IF NOT EXISTS tracks (
	track_id text PRIMARY KEY,
	track_name text NOT NULL,
	track_uri text NOT NULL,
	FOREIGN KEY(album_id) REFERENCES albums(album_id)
);