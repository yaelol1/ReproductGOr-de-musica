package main

import (
	// "fmts "
	"strconv"
	"errors"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Database *sql.DB
	Path string
}

// CreateDatabase creates a database with the tables and rows needed.
func NewDatabase (path string) *Database{
	database, _ := sql.Open("sqlite3", path)

	// execute
	tables := ` CREATE TABLE types ( id_type INTEGER PRIMARY KEY , description TEXT );`
	create, _ := database.Prepare(tables)
	create.Exec()

	tables = ` INSERT INTO types VALUES (0 , "Person");`
	create, _ = database.Prepare(tables)
	create.Exec()


	tables = ` INSERT INTO types VALUES (1 , "Group"); `
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` INSERT INTO types VALUES (2 , "Unknown");`
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` CREATE TABLE performers ( id_performer INTEGER PRIMARY KEY , id_type INTEGER ,
		    name TEXT , FOREIGN KEY ( id_type ) REFERENCES types ( id_type ) );`
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` CREATE TABLE persons ( id_person INTEGER PRIMARY KEY , stage_name TEXT , real_name TEXT ,
		     birth_date TEXT , death_date TEXT );`
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` CREATE TABLE groups ( id_group INTEGER PRIMARY KEY , name TEXT , start_date TEXT , end_date TEXT );`
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` CREATE TABLE albums ( id_album INTEGER PRIMARY KEY , path TEXT , name TEXT , year INTEGER );`
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` CREATE TABLE rolas ( id_rola INTEGER PRIMARY KEY , id_performer INTEGER , id_album INTEGER ,
		    path TEXT , title TEXT , track INTEGER , year INTEGER , genre TEXT , FOREIGN KEY ( id_performer )
		    REFERENCES performers ( id_performer ) , FOREIGN KEY ( id_album ) REFERENCES albums ( id_album ) );`
	create, _ = database.Prepare(tables)
	create.Exec()

	tables = ` CREATE TABLE in_group ( id_person INTEGER , id_group INTEGER , PRIMARY KEY ( id_person , id_group ) ,
		    FOREIGN KEY ( id_person ) REFERENCES persons ( id_person ) , FOREIGN KEY ( id_group ) REFERENCES groups ( id_group ) ); `

	create, _ = database.Prepare(tables)
	create.Exec()

	return &Database{
		Database: database,
		Path: path,
	}
}

// OpenDatabase opens the sqlite3 database with the given path to return a Database struct.
func OpenDatabase (path string) (*Database, error){
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Database{
		Database: database,
		Path: path,
	}, nil
}

// AddSong adds a song to the database, returns true if the song was
// added or false if it was already in the database.
func (database *Database) AddSong(songToAdd *Song) bool{
	if _, err := database.FindSong(songToAdd); err == nil{
		return false
	}

	// TODO: insertar canción
	statement, _ := database.Database.Prepare(`INSERT INTO rolas (id_performer, id_album, path, title, track, year, genre)
	VALUES (?, ?, ?, ?, ?, ?, ?)`)
	statement.Exec()
	return true
}

// FindSong finds a Song in the database and returns the song.
func (database *Database) FindSong(songToSearch *Song) (*Song, error){
	url := songToSearch.path
	rows, _ := database.Database.Query(" SELECT id_rola, id_performer, id_album, path, title, year, genre FROM rolas WHERE path = '"+ url + "'")

	var id_rola, id_performer, id_album, year int
	var path, title, genre string
	hadNext := false

	for rows.Next(){
		rows.Scan(&id_rola, &id_performer, &id_album, &path, &title, &year, &genre)
		log.Printf("id: %v, %v, %v, path: %v, title: %v, year: %v, genre: %v", id_rola, id_performer, id_album, path, title, year, genre )
		hadNext = true
	}

	if !hadNext {
		return nil, errors.New("Song not found.")
	}

	return &Song{}, nil
}


// addPerformer adds a performer if it doesn't exist.
func (database *Database) findPerformer(name string) (int, error) {
	return 0, nil

}

// addAlbum adds an album to the database and returns an Album and an error.
func (database *Database) findAlbum(album *Album) (*Album, error) {
	if ( album.id_Album != 0 ) {

		id := strconv.Itoa(album.id_Album)
		rows, _ := database.Database.Query("SELECT  path, name, year FROM albums WHERE id ="+id)

		var  year int
		var path, name string

		for rows.Next(){
			rows.Scan(&path, &name, &year)
		}

		return &Album{
			id_Album: album.id_album,
			path: path,
			name: name,
			year: year,
		}, nil
	}


	return nil, nil
}
