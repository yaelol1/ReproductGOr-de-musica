package main

import (
	"testing"
	// "log"
	"reflect"
	"os"
)

func TestNewDatabase (t *testing.T) {
	dirname,_  := os.UserConfigDir()
	dirname  += "/reproductgor"
	os.Mkdir( dirname, 0750 )

	defer os.Remove( dirname+"/music_test.db" )

	NewDatabase( dirname+"/music_test.db" )
	if _, err := os.Stat(dirname+"/music_test.db"); err != nil {
		t.Fatalf("Database was not created at %v", dirname+"/music_test.db")
	}
}

func TestOpenDatabase (t *testing.T) {
	_, databasePath := testingDatabase()
	defer func() {
		os.Remove( databasePath )

		if r := recover(); r != nil {
			t.Fatalf("Database was not opened %v", r)
		}
	}()

	OpenDatabase( databasePath )
}

func TestAddSong (t *testing.T){
	_, databasePath := testingDatabase()

	defer os.Remove( databasePath )

	database, _ := OpenDatabase( databasePath )
	tests := []Song{
		{
			performers: "Elthon Jon",
			album: "elthon jon Album",
			path: "~/Music/song.mp3",
			title: "some song",
			year: 1983,
			genre: "pop",
		},
		{
			performers: "Kero Kero Bonito",
			album: "Bonito Generation",
			path: "~/Music/songBonito.mp3",
			title: "some song",
			year: 1983,
			genre: "pop",
		},
	}

	for _, tt := range tests {
		database.AddSong(&tt)
		if got, _ := database.FindSong(&tt); !reflect.DeepEqual(got, &tt) {
			t.Errorf("AddSong(%v) = %v, want %v", tt, got, tt)
		}

	}
}

func testingDatabase() (dirname string, databasePath string){
	dirname, _  = os.UserConfigDir()
	dirname  += "/reproductgor"
	os.Mkdir( dirname, 0750 )
	databasePath = dirname + "/music_test.db"

	NewDatabase( databasePath )
	return dirname, databasePath
}
