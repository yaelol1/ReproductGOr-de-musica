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

	log.Printf("DEBUG: addAlbum enter")
	album, _ = database.addAlbum(album)
	id_album := album.id_album
	log.Printf("DEBUG: addAlbum exit")

	performer := &Performer{
		name: songToAdd.performers,
		id_type: 3,
	}

	performer, _ = database.addPerformer(performer)
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
func (database *Database) addPerformer(performer *Performer) (*Performer, error) {
	if performerFound, err := database.findPerformer(performer); err == nil{
		return performerFound, nil
	}

	if !performer.Addable() {
		return nil, errors.New("The performer does not have enough data to be added")
	}

	statement, _ := database.Database.Prepare(`INSERT INTO performers (id_type, name)
	VALUES (?, ?)`)
	statement.Exec(performer.id_type, performer.name)

	performer, _ = database.findPerformer(performer)

	return performer, nil
}

// findPerformer adds a performer if it doesn't exist.
func (database *Database) findPerformer(performer *Performer) (*Performer, error) {
	if performer.id_performer != 0  {

		id := strconv.Itoa(performer.id_performer)
		rows, _ := database.Database.Query("SELECT id_type, name FROM albums WHERE id ="+id)

		var id_type int
		var name string

		for rows.Next(){
			rows.Scan(&id_type, &name)
		}

		performer.id_type = id_type
		performer.name = name

		return performer, nil
	}

	// The Performer doesn't have an id, nor a name, therefore cannot be searched.
	if performer.name == "" {
		return nil, errors.New("The Performer must have a Name or an id to be found in the database.")
	}

	// The performer will be searched by name
	var id_performer, id_type int
	hadNext := false

	rows, _ := database.Database.Query("SELECT id_performer, id_type FROM performers WHERE name ="+performer.name)

	for rows.Next() {
		rows.Scan(&id_performer, &id_type)
		hadNext = true
	}

	// The performer wasn't found in the database.
	if !hadNext {
		return nil, errors.New("Performer not found.")
	}

	performer.id_type = id_type
	performer.id_performer = id_performer

	return performer, nil
}


// addAlbum adds an album if it doesnt' exist
func (database *Database) addAlbum(album *Album) (*Album, error) {

	log.Printf("DEBUG: findAlbum enter")
	if albumFound, err := database.findAlbum(album); err == nil{

		log.Printf("DEBUG: findAlbum inside")
		return albumFound, errors.New("Album already in database")
	}
	log.Printf("DEBUG: findAlbum exit")

	if !album.Addable() {
		return nil, errors.New("Album needs a name, path and yaer to be added")
	}


	statement, _ := database.Database.Prepare(`INSERT INTO albums (name, path, year)
	VALUES (?, ?, ?)`)
	statement.Exec(album.name, album.path, album.year)

	album, _ = database.findAlbum(album)

	return album, nil
}

// findAlbum adds an album to the database and returns the album given with the database entries Album and an error.
func (database *Database) findAlbum(album *Album) (*Album, error) {
	// The album given has an id.
	log.Printf("DEBUG: findAlbum if")
	if album.id_album != 0 && album.year != 0 {

		log.Printf("DEBUG: findAlbum if -> true")
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
	log.Printf("DEBUG: findAlbum if -> else")

	// The Album doesn't have an id, nor a name, therefore cannot be searched.
	if album.name == "" {
		return nil, errors.New("The album must have a Name or an id to be found in the database.")
	}

	// The album will be searched by name
	var  id_album, year int
	var path string
	hadNext := false

	log.Printf("DEBUG: findAlbum query")
	rows, err := database.Database.Query("SELECT id_album, path, year FROM albums WHERE name ="+album.name)

	log.Printf("DEBUG: findAlbum for rows.Next inicia")
	for rows.Next() {
		log.Printf("DEBUG: findAlbum for rows.Next")
		errRows := rows.Scan(&id_album, &path, &year)
		hadNext = true
	}

	// The album wasn't found in the database.
	if !hadNext {
		log.Printf("DEBUG: findAlbum if !hadNext -> true")
		return nil, errors.New("Album not found.")
	}

	log.Printf("DEBUG: findAlbum album =")
	album.id_album = id_album
	album.path = path
	album.year = year

	return album, nil
}
