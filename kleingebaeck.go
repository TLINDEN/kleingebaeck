package main

var manpage = `
NAME
    kleingebaeck - kleinanzeigen.de backup tool

SYNOPSYS
        Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
        Options:
        --user    -u <uid>      Backup ads from user with uid <uid>.
        --debug   -d            Enable debug output.
        --verbose -v            Enable verbose output.
        --outdir  -o <dir>      Set output dir (default: current directory)
        --limit   -l <num>      Limit the ads to download to <num>, default: load all.
        --config  -c <file>     Use config file <file> (default: ~/.kleingebaeck).
        --manual  -m            Show manual.
        --help    -h            Show usage.

DESCRIPTION
    This tool can be used to backup ads on the german ad page
    <https://kleinanzeigen.de>.

    It downloads all (or only the specified ones) ads of one user into a
    directory, each ad into its own subdirectory. The backup will contain a
    textfile Adlisting.txt which contains the ad contents such as title,
    body, price etc. All images will be downloaded as well.

CONFIGURATION
    You can create a config file to save typing. By default
    "~/.kleingebaeck" is being used but you can specify one with "-c" as
    well. We use TOML as our configuration language. See
    <https://toml.io/en/>.

    Format is pretty simple:

        user = 1010101
        loglevel = verbose
        outdir = "test"
        template = """
        Title: {{.Title}}
        Price: {{.Price}}
        Id: {{.Id}}
        Category: {{.Category}}
        Condition: {{.Condition}}
        Created: {{.Created}}

        {{.Text}}
        """

    Be carefull if you want to change the template. The variable is a
    multiline string surrounded by three double quotes. You can left out
    certain fields and use any formatting you like. Refer to
    <https://pkg.go.dev/text/template> for details how to write a template.

    If you're on windows and want to customize the output directory, put it
    into single quotes to avoid the backslashes interpreted as escape chars
    like this:

        outdir = 'C:\Data\Ads'

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
