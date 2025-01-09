

fix meta data 

this program reads MP3 and flac files
and writes out commands to fix(remove or replace) unicode characters in the
artist, title and album metadata tags

Read the meta data and find non-ASCII characters (Unicode higher than 0x7F)
Output commands to fix the tags to just use printable ASCII

this is not multi-processing safe
this is not general, it uses a specific look up table rather
than the official "PRECIS"
(Preparation, Enforcement, and Comparison of Internationalized Strings in Application Protocols)
PRECIS is documented in RFC7564.

for example Joshua Bell's music often uses Cyrillic or Polish, which this does not handle.


Changes ’ (U+2019) (curved single quote) to straight quote
U+2010 hyphen to minus sign
… (U+2026) three dots to nothing
U+2013 	– 	En dash 	0903
U+2014 	— 	Em dash 	0904
U+2015 	― 	Horizontal bar

