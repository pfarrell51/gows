// Package filenametool is useful for files that do not contain accurate
// mp3/Vorbis/Flac metadata tags. It reads the files name and path and
// attempts to figure out what artist, album and song title is
// contained in the file.
//
// The basic operation is to run the executable, give it a switch/flag to say what
// you want to do, and specify a folder/directory that will be processed
//
// Usage of mm/mm: [flags] directory-spec
//
//  -a    artist map -  list artist not in source code. The list of artists is
//			contained in a normal text file named data/artists.txt 
//  -de   debug on, output lots of error messages that help debugging
//  -dup   duplicate song.  attempts to identify duplicate songs, does not work well 
//  -h    help
//  -j    output metadata as json
//  -l    list - list files
//  -n    nogroup - list files that do not have an artist/group in the title
//  -r    rename - output rename from parsed file name
//  -z    list artist names one per line
//
// default is to list files that need love.
//
// none of these operations change the files being processed, they are all read-only
// the -r command creates a command that can be used to rename the files
// into the new format.
package filenametool
