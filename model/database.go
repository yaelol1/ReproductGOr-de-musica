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

	album := &Album{
		name: songToAdd.album,
		path: songToAdd.path,
	}
	database.addAlbum(album)
	id_album := album.id_album

	performer, _ := database.addPerformer(songToAdd.performers)
	id_performer := performer.id_performer

	statement.Exec(id_performer, id_album, songToAdd.path, songToAdd.title, 0, songToAdd.year, songToAdd.genre)
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
func (database *Database) addPerformer(name string) (int, error) {
	return 0, nil
}

// findPerformer adds a performer if it doesn't exist.
func (database *Database) findPerformer(name string) (int, error) {
	return 0, nil
}


// addAlbum adds an album if it doesnt' exist
func (database *Database) addAlbum(album *Album) (*Album, error) {
	if albumFound, err := database.findAlbum(album); err == nil{
		return albumFound, errors.New("Album already in database")
	}

	if(!album.Addable()){
		return nil, errors.New("Album needs a name, path and yaer to be added")
	}


	statement, _ := database.Database.Prepare(`INSERT INTO albums (name, path, year)
	VALUES (?, ?, ?)`)
	statement.Exec(album.name, album.path, album.year)

	findAlbum(album)

	return album, nil
}

// findAlbum adds an album to the database and returns the album given with the database entries Album and an error.
func (database *Database) findAlbum(album *Album) (*Album, error) {
	// The album given has an id.
	if ( album.id_album != 0 && album.year != 0 ) {

		id := strconv.Itoa(album.id_album)
		rows, _ := database.Database.Query("SELECT  path, name, year FROM albums WHERE id ="+id)

		var  year int
		var path, name string

		for rows.Next(){
			rows.Scan(&path, &name, &year)
		}

		album.path = path
		album.name = name
		album.year = year

		return album, nil
	}

	// The Album doesn't have an id, nor a name, therefore cannot be searched.
	if(album.name == ""){
		return nil, errors.New("The album must have a Name or an id to be found in the database.")
	}

	// The album will be searched by name
	var  id_album, year int
	var path string
	hadNext := false

	rows, _ := database.Database.Query("SELECT id_album, path, year FROM albums WHERE name ="+album.name)

	for rows.Next(){
		rows.Scan(&id_album, &path, &year)
		hadNext = true
	}

	album.id_album = id_album
	album.path = path
	album.year = year

	// The album wasn't found in the database.
	if(!hadNext){
		return nil, errors.New("Album not found.")
	}

	return album, nil
}
