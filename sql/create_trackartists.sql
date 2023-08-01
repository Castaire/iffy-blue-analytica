CREATE TABLE IF NOT EXISTS trackArtists (
    id integer PRIMARY KEY AUTOINCREMENT,
    FOREIGN KEY(track_id) REFERENCES tracks(track_id),
    FOREIGN KEY(artist_id) REFERENCES artists(artist_id)
);