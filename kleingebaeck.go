package main

var manpage = `
NAME
    kleingebaeck - kleinanzeigen.de backup tool

SYNOPSYS
        This is kleingebaeck, the kleinanzeigen.de backup tool.
        Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
        Options:
        --user,-u <uid>        Backup ads from user with uid <uid>.
        --debug, -d            Enable debug output.
        --verbose,-v           Enable verbose output.
        --output-dir,-o <dir>  Set output dir (default: current directory)
        --manual,-m            Show manual.
        --config,-c <file>     Use config file <file> (default: ~/.kleingebaeck).

DESCRIPTION
    This tool can be used to backup ads on the german ad page
    <https://kleinanzeigen.de>.

    It downloads all (or only the specified ones) ads of one user into a
    directory, each ad into its own subdirectory. The backup will contain a
    textfile Adlisting.txt which contains the ad contents such as title,
    body, price etc. All images will be downloaded as well.

CONFIGURATION
    You can create a config file to save typing. By default
    "~/.kleingebaeck.hcl" is being used but you can specify one with "-c" as
    well.

    Format is simple:

        user = 1010101
        verbose = true
        outdir = "test"

SETUP
    To setup the tool, you need to lookup your userid on kleinanzeigen.de.
    Go to your ad overview page while NOT being logged in:

        https://www.kleinanzeigen.de/s-bestandsliste.html?userId=XXXXXX

    The XXXXX part is your userid.

    Put it into the configfile as outlined above. Also specify an output
    directory. Then just execute "kleingebaeck".

    You can use the -v option to get verbose output or -d to enable
    debugging.

BUGS
    In order to report a bug, unexpected behavior, feature requests or to
    submit a patch, please open an issue on github:
    <https://github.com/TLINDEN/kleingebaeck/issues>.

    Please repeat the failing command with debugging enabled "-d" and
    include the output in the issue.

LIMITATIONS
    The "kleingebaeck" doesn't currently check if it has downloaded a file
    already, so it downloads everything again every time you execute it. Be
    aware of it. This will change in the future.

    Also there's currently no parallelization implemented. This will change
    in the future.

LICENSE
    Licensed under the GNU GENERAL PUBLIC LICENSE version 3.

Author
    T.v.Dein <tom AT vondein DOT org>

`
