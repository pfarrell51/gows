# music/tagtool
music related programs, utilites, etc. for song files that contain "tags"
or metadata tags in the mp3/Vorbis/FLAC formats.
These are mostly about managing a library of recorded music.
So far, these are useful for music "ripped" and "tagged" from
commercial audio CDs.

This program assume that the files have been tagged by a utility such
as Musicbrainz's Picard (https://picard.musicbrainz.org/)

 Pat has about 1000 CDs collected from 1985. He still occasionally
 buys new ones, but most of his music is from last century

 The basic operation is to run the executable, give it a switch/flag to say what you want to do, and specify a folder/directory that will be processed.

 as usual, there is a package "tagtool" of software that performs the operations and a subdirectory that contains the main.go file that is the executable. The subdirectory for tagtools is named "mm" so the executable that is built by "go build" will be called "mm" or "mm.exe"

 Usage of mm/mm: [flags] directory-spec

  -a    artist map -  list artist not known in source code. The known 
  artists are kept in a text file in the data subdirectory, named data/artists.txt

  -c    Album track order - output cp command to copy songs in track order so that you can play the songs in the same order as they were in an album. Many albums in the late 60s and 70s had a specific flow for the songs


 -csv       output CSV format

  -de      debug on

  -dup     locate duplicate song based on title, album & Artist on

  -duptitle     duplicate song based on just the title on

  -h    help

  -i    inventory - basic inventory (handy with -i -csv -sn)

  -j    output metadata as json

  -l    list - list files

  -ng      nogroup - list files that do not have an artist/group in the title

  -nt      notags - list files that do not have any meta tags

  -r    rename - output rename from internal metadata

  -s    summary - print summary statistics

  -sn     show no song titles (in inventory and other listings)

  -z    list artist names one per line

default is to list files that need love.

none of these operations change the files being processed, they are all read-only.  The -r command creates a command that can be used to rename the files into the new format.

You can have multiple switches (where it makes sense) in a single use.
For example,
     mm/mm -i  directory-spec
will list the artist, album and song title for each songs file (.mp3, .flac) in
the tree of directories

     mm/mm -i -csv directory-spec
will do the same, outputting the list in .csv (comma separated values) so the resulting file can easily be read into a spreadsheet or database

     mm/mm -i -csv -sn directory-spec
will show no song titles, i.e. just the artist and album name.
This is handy when you want to compare you inventory to your physical collection of CDs in a box somewhere

