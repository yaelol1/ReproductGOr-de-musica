// package model mines mp3 files from a given folder and adds them to a database,
// also has the libraries to play the mp3 files.
package main

import (
	// "fmt"
	// "errors"
	"strconv"
	"github.com/bogem/id3v2/v2"
	"log"
	"path/filepath"
	"os"
)

// A Song is a file that contains the basic info of the song.
type Song struct {
	id_song int
	performers string
	album string
	path string
	title string
	year int
	genre string
}

// An Album is a collection of songs
type Album struct {
	id_Album int
	name string
	path string
	year int
}

// TODO: performer or a band and a artist
// type Performer stuct {
	// id_performer int
// }


// mine walks recursively the path given, to find every .mp3 song, to store it
// in the database given.
func Mine(path string, database *Database) {
	filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			newSong, err := NewSong(path, info)

			// The file didn't have any tags so it wont be stored
			if err != nil {
				log.Print( err )
				return err
			}
			// TODO: Does it affect?
			go database.AddSong( newSong )
		}
		return nil
	})
}

// NewSong takes a file and its path to create a Song Struct and return it.
func NewSong(path string, info os.FileInfo) (*Song, error){

	// opening id3 tag
	tag, err := id3v2.Open( path , id3v2.Options{Parse: true} )
	if err != nil {
		return nil, err
	}

	defer tag.Close()

	year, _ := strconv.Atoi( tag.Year() )

	return &Song{
		performers: tag.Artist(),
		album: tag.Album(),
		path: path,
		title: tag.Title(),
		year: year,
		genre: tag.Genre(),
	}, nil
}
