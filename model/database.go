// package model mines mp3 files from a given folder and adds them to a database, also has the libraries to play the mp3 files.
package reproductgor/model

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Database *sql.DB
}

// CreateDatabase creates a database with the tables and rows needed.
func NewDatabase (path string){

}
