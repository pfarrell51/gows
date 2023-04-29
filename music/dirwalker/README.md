# music utility dirwalker
the utility program walks music file directories and execute a "verb" against each music file.
(This could be done with a proper "find" command, but that syntax is so strange.....)

The basic operation is to run the executable, give it a "verb" to say what you want to do, and specify a folder/directory that will be processed.

Usage of dw/dw [verb] indirectory-spec outdirecgtory-spec extension-spec

Supported verbs are

  ffmpeg	convert a flac file to mp3

  sox       use the sox "compand" command to compress the dynamic range of the audio

  both      do to commands, ffmpeg and then sox

none of these operations change the files being processed, they are all read-only.  The execution of a verb 
creates a command that can be used to rename the files into the new format.


