# gopro rename goprorename or GoPro rename

Walk the argument directory path looking for .mp4 files generated
by a GoPro camera, creating commands to rename the files
from the GoPro naming convention to names that make sense to humans
and to utilities such as RaceRender(tm).

The basic operation is to run the executable, and specify a folder/directory that will be processed.
For example if the executable is called 'goprorename' then one can run the program
and specify a directory containing MP4 files generated by a GoPro camera:

   goprorename  /Videos/GoPro9/April10

The program will write out commands to rename the files.
Example output will look like:  
    # GX010391.MP4 GX010391.MP4  
    # GX020391.MP4 GX020391.MP4  
    # GX030391.MP4 GX030391.MP4  
    # GX040391.MP4 GX040391.MP4  
    mv GX010392.MP4 GX050391.MP4  
    mv GX020392.MP4 GX060391.MP4  
    mv GX010393.MP4 GX070391.MP4  
    mv GX020393.MP4 GX080391.MP4  
    mv GX010394.MP4 GX090391.MP4  

Executing the program will not change the files being processed, they are all read-only.
The output is a set of commands that will rename the GoPro files.


When you want load the renamed videos into RaceRender, you only need to enter the first
file name in the "Add" function. RaceRender will notice the additional files and
ask if you want to load them too. Click "Yes"

Internally, it creates  a map of old and new filenames
then calling the processMap function. The output are shell/cmd commands
to rename the files in a human sensible matter. This really helps
managing the files. Plus it makes loading them into a utility such as
RaceRender a lot more pleasant.


As usual, there is a package "goprorename" of software that performs the operations
and a subdirectory that contains the main.go file that is the executable.
the subdirectory for tagtools is named "mg" so the executable that is built
by "go build" will be called "mg" or "mg.exe"



Limitation:
The initial version of this utility can only handle 99 files in the directory.
It would be easy to change this but 99 has been more than enough for
my uses. If someone needs more, let me know.

