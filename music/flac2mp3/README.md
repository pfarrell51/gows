# music utility dirwalker
the utility program walks music file directories and execute a "verb" against each music file.
(This could be done with a proper "find" command, but that syntax is so strange.....)

The basic operation is to run the executable, give it a "verb" to say what you want to do, 
and specify a folder/directory that will be processed.

Usage:
     dw/dw --flags indirectory-spec outdirecgtory-spec extension-spec


Alternative usage of:
        dw/dw  indirectory-spec mp3 

The execution of the program creates a command that can be used to rename the files into the new format.

