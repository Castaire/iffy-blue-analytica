package db

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const DBFolder = "../dbs"
const SQLFolder = "sql"

func CreateTables(dropAll bool) error {

	tableNames := [4]string{"create_albums.sql", "create_artists.sql", "create_tracks.sql", "create_trackartists.sql"}
	dbNames := [4]string{"albums.db", "artists.db", "tracks.db", "trackartists.db"}

	for i := 0; i < len(tableNames); i++ {
		bCmd, err := os.ReadFile(path.Join(SQLFolder, tableNames[i]))
		if err != nil {return err}

		if dropAll {
			// drop before re-creating table
			err = dropTable(tableNames[i], dbNames[i])
			if err != nil {return err}
		}

		err = dbExec(path.Join(DBFolder, dbNames[i]), string(bCmd))
		if err != nil {return err}
	}

	return nil
}

func dropTable(tableName string, dbName string) error {
	cmd := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	return dbExec(path.Join(DBFolder, dbName), cmd)
}

// RETURN: (last insert ID, error)
func InsertTracks(tracks []Track) (int64, error) {

	sqlStr := "INSERT INTO tracks(track_id, track_name, track_uri, album_id) VALUES "
	vals := []interface{}{}

	for _, track := range tracks {
		sqlStr += "(?, ?, ?), "
		vals = append(vals, track.TrackID, track.Name, track.URI, track.AlbumID)
	}

	sqlStr = strings.TrimRight(sqlStr, ", ")
	return dbExecPrepared(path.Join(DBFolder, "tracks.db"), sqlStr, vals)
}

func dbExecPrepared(dbPath string, sqlStr string, vals []interface{}) (int64, error) {
	
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {return 0, err}
	defer db.Close()

	stmt, err := db.Prepare(sqlStr)
	if err != nil {return 0, err}

	res, err := stmt.Exec(vals...)
	if err != nil {return 0, err}

	lastInsertId, err := res.LastInsertId()
	if err != nil {return 0, err}
	return lastInsertId, err
}

func dbExec(dbPath string, cmd string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {return err}
	defer db.Close()

	_, err = db.Exec(cmd)
	return err
}