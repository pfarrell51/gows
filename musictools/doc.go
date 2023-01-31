//Package musictools is Pat's tools to organize and manage his music library
// Pat has about 1000 CDs collected from 1985. He still occasionally
// buys new ones, but most of his music is from last century
//
// The basic operation is to run the executable, give it a switch/flag to say what 
// you want to do, and specify a folder/directory that will be processed
// 
// Usage of mm/mm: [flags] directory-spec
//  -a    artist map -  list artist not in source code (gpmap)
//  -d    debug on
//  -j    output metadata as json
//  -l    list - list files
//  -n    nogroup - list files that do not have an artist/group in the title
//  -r    rename - output command to perform rename function on needed files
//  -z    list artist names one per line
//default is to list files that need love.
//
// none of these operations change the files being processed, they are all read-only
// the -r command creates a command that can be used to rename the files
// into the new format.
package musictools

