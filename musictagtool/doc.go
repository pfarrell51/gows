// Package musictools is Pat's tools to organize and manage his music library
// Pat has about 1000 CDs collected from 1985. He still occasionally
// buys new ones, but most of his music is from last century
//
// The basic operation is to run the executable, give it a switch/flag to say what
// you want to do, and specify a folder/directory that will be processed
//
// Usage of mm/mm: [flags] directory-spec
//
//		-a    artist map -  list artist not in source code (gpmap)
//		-de    debug on
//	 	-dup  try to identify duplicate songs. Note, the same song title by two different
//	       artists does not count as a duplicate
//	 	-fr   rename - output command to perform rename function on needed files, the -mr is
//	         much more reliable, use Musicbrainz' Picard utility to set the metadata
//		-j    output metadata as json
//		-l    list - list files
//		-mr   rename based on Metadata. output commands to rename files based on mp3 metadata
//		-n    nogroup - list files that do not have an artist/group in the title
//		-z    list artist names one per line
//
// default is to list files that need love.
//
// none of these operations change the files being processed, they are all read-only
// the -r command creates a command that can be used to rename the files
// into the new format.
package musictagtool
