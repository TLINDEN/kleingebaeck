package main

var manpage = `
NAME
    kleingebaeck - kleinanzeigen.de backup tool

SYNOPSYS
        Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
        Options:
        -u --user    <uid>      Backup ads from user with uid <uid>.
        -d --debug              Enable debug output.
        -v --verbose            Enable verbose output.
        -o --outdir  <dir>      Set output dir (default: current directory)
        -l --limit   <num>      Limit the ads to download to <num>, default: load all.
        -c --config  <file>     Use config file <file> (default: ~/.kleingebaeck).
           --ignoreerrors       Ignore HTTP errors, may lead to incomplete ad backup.
        -f --force              Overwrite images and ads even if the already exist.
        -m --manual             Show manual.
        -h --help               Show usage.
        -V --version            Show program version.

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
        useragent = "Mozilla/5.0"
        template = """
        Title: {{.Title}}
        Price: {{.Price}}
        Id: {{.ID}}
        Category: {{.Category}}
        Condition: {{.Condition}}
        Created: {{.Created}}

        {{.Text}}
        """

    Be careful if you want to change the template. The variable is a
    multiline string surrounded by three double quotes. You can left out
    certain fields and use any formatting you like. Refer to
    <https://pkg.go.dev/text/template> for details how to write a template.
    Also read the TEMPLATES section below.

    If you're on windows and want to customize the output directory, put it
    into single quotes to avoid the backslashes interpreted as escape chars
    like this:

        outdir = 'C:\Data\Ads'

TEMPLATES
    Various parts of the configuration can be modified using templates: the
    output directory, the ad directory and the ad listing itself.

  OUTPUT DIR TEMPLATE
    The config varialbe "outdir" or the command line parameter "-o" take a
    template which may contain:

    "{{.Year}}"
    "{{.Month}}"
    "{{.Day}}"

    That way you can create a new output directory for every backup run. For
    example:

        outdir = "/home/backups/ads-{{.Year}}-{{.Month}}-{{.Day}}"

    Or using the command line flag:

        -o "/home/backups/ads-{{.Year}}-{{.Month}}-{{.Day}}"

    The default value is "." - the current directory.

  AD DIRECTORY TEMPLATE
    The ad directory name can be modified using the following ad values:

    {{.Price}}
    {{.ID}}
    {{.Category}}
    {{.Condition}}
    {{.Created}}
    {{.Slug}}
    {{.Text}}

    It can only be configured in the config file. By default only
    "{{.Slug}}" is being used, this is the title of the ad in url format.

  AD NAME TEMPLATE
    The name of the directory per ad can be tuned as well:

    "{{.Year}}"
    "{{.Month}}"
    "{{.Day}}"
    "{{.Slug}}"
    "{{.Category}}"
    "{{.ID}}"

  AD TEMPLATE
    The ad listing itself can be modified as well, using the same variables
    as the ad name template above.

    This is the default template:

        Title: {{.Title}}
        Price: {{.Price}}
        Id: {{.ID}}
        Category: {{.Category}}
        Condition: {{.Condition}}
        Type: {{.Type}}
        Created: {{.Created}}
        Expire: {{.Expire}}
    
        {{.Text}}

    The config parameter to modify is "template". See example.conf in the
    source repository. Please take care, since this is a multiline string.
    This is how it shall look if you modify it:

        template="""
        Title: {{.Title}}
    
        {{.Text}}
        """

    That is, the content between the two """ chars is the template.

SETUP
    To setup the tool, you need to lookup your userid on kleinanzeigen.de.
    Go to your ad overview page while NOT being logged in:

        https://www.kleinanzeigen.de/s-bestandsliste.html?userId=XXXXXX

    The XXXXX part is your userid.

    Put it into the configfile as outlined above. Also specify an output
    directory. Then just execute "kleingebaeck".

    You can use the -v option to get verbose output or -d to enable
    debugging.

ENVIRONMENT VARIABLES
    The following environment variables are considered:

        KLEINGEBAECK_USER
        KLEINGEBAECK_DEBUG
        KLEINGEBAECK_VERBOSE
        KLEINGEBAECK_OUTDIR
        KLEINGEBAECK_LIMIT
        KLEINGEBAECK_CONFIG
        KLEINGEBAECK_IGNOREERRORS

    Please note, that they take precedence over config file, but commandline
    flags take precedence over env!

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
    Copyright 2023-2025 Thomas von Dein

    This program is free software: you can redistribute it and/or modify it
    under the terms of the GNU General Public License as published by the
    Free Software Foundation, either version 3 of the License, or (at your
    option) any later version.

    This program is distributed in the hope that it will be useful, but
    WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
    Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program. If not, see <http://www.gnu.org/licenses/>.

Author
    T.v.Dein <tom AT vondein DOT org>

`
