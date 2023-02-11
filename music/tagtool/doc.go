// Package tagtools is a tool that uses standard MP3/Vorbis/Flac meta data tags
// to manage music files. These file are typically "ripped" from audio CDs.
// Pat has about 1000 CDs collected from 1985. He still occasionally
// buys new ones, but most of his music is from last century
//
// The basic operation is to run the executable, give it a switch/flag to say what
// you want to do, and specify a folder/directory that will be processed.
// 
// as usual, there is a package "tagtool" of software that performs the operations
// and a subdirectory that contains the main.go file that is the executable.
// the subdirectory for tagtools is named "mm" so the executable that is built
// by "go build" will be called "mm" or "mm.exe"
//
// Usage of mm/mm: [flags] directory-spec
//
//  -a    artist map -  list artist not known in source code. The known artists are
//			kept in a text file in the data subdirectory, named artists.txt 
//  -c    Album track order - output cp command to copy songs in track order so that
//			you can play the songs in the same order as they were in an album.
//			Many albums in the late 60s and 70s had a specific flow for the songs,
//			they were not just a collection of singles.
//  -de   debug on, turns on all sorts of debugging messages
//  -dup  duplicate attempts to identify duplicate songs
//  -h    help
//  -j    output metadata as json
//  -l    list - list files
//  -ng   nogroup - list files that do not have an artist/group in the title
//  -nt   notags - list files that do not have any meta tags
//  -re   rename - output rename from internal metadata. This will generate 
//			rename commands for the files so that the file name reflects the song title and artist
//			from the metadata tags
//  -z    list artist names one per line
//

// default is to list files that need love.
//
// none of these operations change the files being processed, they are all read-only
// the -r command creates a command that can be used to rename the files
// into the new format.
package tagtool
