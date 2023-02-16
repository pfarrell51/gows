# gopro rename goprorename or GoPro rename

Walk the argument directory path looking for .mp4 files generated
by a GoPro camera, creating a map of old and new filenames
then calling the processMap function. The output are shell/cmd commands
to rename the files in a human sensible matter. This really helps
managing the files. Plus it makes loading them into a utility such as
RaceRender a lot more pleasant.


The basic operation is to run the executable, and specify a folder/directory that will be processed.
 
As usual, there is a package "goprorename" of software that performs the operations
and a subdirectory that contains the main.go file that is the executable.
the subdirectory for tagtools is named "mg" so the executable that is built
by "go build" will be called "mg" or "mg.exe"


None of these operations change the files being processed, they are all read-only.
The output is a set of commands that will rename the GoPro files

When you want load the renamed videos into RaceRender, you only need to enter the first
file name in the "Add" function. RaceRender will notice the additional files and
ask if you want to load them too. Click "Yes"

