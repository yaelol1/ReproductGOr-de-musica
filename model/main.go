// package model mines mp3 files from a given folder and adds them to a database, also has the libraries to play the mp3 files.
package main

import (
	"fmt"
	"log"
	"os"
)

// main tests the database.
func main(){
	defer recoverInvalidAccess()

	dirname, err := os.UserConfigDir()
	dirname += "/reproductgor"
	if err != nil {
		log.Print( err )
	}

	errMk := os.Mkdir(dirname, 0750)
	if errMk != nil {
		log.Print( err )
	}

	// checks if the database exists
	if _, err = os.Stat(dirname+"/music.db"); err != nil {
		log.Printf("Database doesn't exists");
		NewDatabase( dirname+"/music.db" )
		log.Print( dirname+"/music.db:  Database created" )
	}

}

// recoverInvalidAccess recovers from an error.
func recoverInvalidAccess() {
	if r := recover(); r != nil {
		fmt.Println( "Recovered", r )
	}
}
