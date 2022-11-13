// package model mines mp3 files from a given folder and adds them to a database, also has the libraries to play the mp3 files.
package main

import (
	// "fmt"
	"log"
	"os"
)

// main tests the database.
func main(){
	defer recoverInvalidAccess()

	// Get dir path
	dirname, err := os.UserConfigDir()
	dirname += "/reproductgor"
	if err != nil {
		log.Print( err )
	}

	// Make Dir
	errMk := os.Mkdir(dirname, 0750)
	if errMk != nil {
		log.Print( err )
	}

	databasePath := dirname+"/music.db"
	var database *Database

	// checks if the database exists
	if _, err = os.Stat(databasePath); err != nil {
		log.Printf("Database doesn't exists");
		NewDatabase( databasePath )
		log.Printf("%v:  Database created", databasePath)
	}

	log.Printf("DEBUG: opening database")
	// Open Database and mine Songs
	database, _ = OpenDatabase(databasePath)

	log.Printf("DEBUG: mining")
	Mine("/home/y421/Music", database)

}

// recoverInvalidAccess recovers from an error.
func recoverInvalidAccess() {
	if r := recover(); r != nil {
		log.Println( "Recovered", r )
	}
}
