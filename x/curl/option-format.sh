
doc<<<EOM
Parse the options into JSON objects, like: [{"short": "-a", "long": "--all", help="do not ignore entries starting with ."}, ...]


LS(1)                                                                                                            User Commands                                                                                                           LS(1)

NAME
       ls - list directory contents

SYNOPSIS
       ls [OPTION]... [FILE]...

DESCRIPTION
       List information about the FILEs (the current directory by default).  Sort entries alphabetically if none of -cftuvSUX nor --sort is specified.

       Mandatory arguments to long options are mandatory for short options too.

       -a, --all
              do not ignore entries starting with .

       -A, --almost-all
              do not list implied . and ..

       --author
              with -l, print the author of each file

       -b, --escape
              print C-style escapes for nongraphic characters

       --block-size=SIZE
              with -l, scale sizes by SIZE when printing them; e.g., '--block-size=M'; see SIZE format below

       -B, --ignore-backups
              do not list implied entries ending with ~

       -c     with -lt: sort by, and show, ctime (time of last change of file status information); with -l: show ctime and sort by name; otherwise: sort by ctime, newest first

       -C     list entries by columns

       --color[=WHEN]
              color the output WHEN; more info below

       -d, --directory
              list directories themselves, not their contents
EOM

data<<<EOB
{
    "model": "llama3.2:latest",
    "prompt": "$doc",
    "format": "json",
    "stream": false,
    "_debug_render_only": true
}
EOB

echo $data

# curl -s localhost:11435/api/generate
