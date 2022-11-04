package main

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/bogem/id3v2/v2"
)

var tests = []Song{
	{
		performers: "Sed diam.",
		album:      "elthon jon Album",
		path:       "~/Music/song.mp3",
		title:      "some song",
		year:       1983,
		genre:      "pop",
	},
	{
		performers: "Kero Kero Bonito",
		album:      "Bonito Generation",
		path:       "~/Music/song.mp3",
		title:      "some song",
		year:       1983,
		genre:      "pop",
	},
	{
		performers: "Nam a sapien.  ",
		album:      "Praesent augue.  ",
		path:       "~/Music/song.mp3",
		title:      "Mauris mollis tincidunt felis.  ",
		year:       1983,
		genre:      "Nulla posuere.  ",
	},
	{
		performers: "Nulla facilisis, risus a rhoncus fermentum,, et dictum nunc justo sit amet elit.  ",
		album:      "Pellentesque condimentum, magna ut suscipitnon luctus diam neque sit amet urna.  ",
		title:      "Suspendisse potenti.  ",
		path:       "~/Music/song.mp3",
		year:       1983,
		genre:      "Donec posuere augue in quam.  ",
	},
	{
		performers: "",
		album:      "",
		path:       "~/Music/song.mp3",
		title:      "",
		year:       1983,
		genre:      "",
	},
}

func TestSongWithTags(t *testing.T) {
	// Creating .mp3 file
	dirname, songPath := testingSong()
	defer os.Remove(songPath)

	// Creating and opening database
	_, databasePath := testingDatabase()
	defer os.Remove(databasePath)
	database, _ := OpenDatabase(databasePath)

	// Changing the tag with the struct and then comparing them
	for _, tt := range tests {
		changeTag(songPath, tt)
		Mine(dirname, database)

		if got, _ := database.FindSong(&tt); !reflect.DeepEqual(got, &tt) {
			t.Errorf("addSong(%v) = %v, want %v", tt, got, tt)
		}
	}
}

// changeTag changes the id3 tag of a given file.
func changeTag(songPath string, newSongTags Song) {

	// opening id3 tag to edit it
	tag, err := id3v2.Open(songPath, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	// Set Tags
	tag.SetArtist(newSongTags.performers)
	tag.SetAlbum(newSongTags.album)
	tag.SetTitle(newSongTags.title)
	year := strconv.Itoa(newSongTags.year)
	tag.SetYear(year)
	tag.SetGenre(newSongTags.genre)

	// Write tag to file.mp3
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}

}

// testingSong creates an .mp3 file in .config/reproductgor to test the file, returns the
// folder and the song path.
func testingSong() (dirname string, songPath string) {
	dirname, _ = os.UserConfigDir()
	dirname += "/reproductgor"
	os.Mkdir(dirname, 0750)
	songPath = dirname + "/song.mp3"

	_, err := os.Create(songPath)

	if err != nil {
		log.Fatal(err)
	}

	return
}
