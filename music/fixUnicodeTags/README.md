fix meta data 
Read the meta data and find non-ASCII characters (Unicode higher than 0x7F)
Output commands to fix the tags to just use printable ASCII

Changes ’ (U+2019) (curved single quote) to straight quote
– (U+2013) en dash to minus sign
‐ (U+2010) hyphen to minus sign
ó (U+00F3) accented o to o
é (U+00E9) accented e to e
ś (U+015B) letter s with acute accent
ù (U+00F9) small u with accent grave
… (U+2026) three dots to nothing
È U+00C8 Capital E accent grave
É (U+00C9) capital accented E to capital E

U+2013 	– 	En dash 	0903
U+2014 	— 	Em dash 	0904
U+2015 	― 	Horizontal bar

