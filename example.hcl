#
# kleingebaeck sample configuration file.
# put this to ~/.kleingebaeck.hcl.
#
# Comments start with the '#' character.

# kleinanzeigen.de user-id. must be an unquoted number
user = 00000000

# enable verbose output (same as -v), may be true or false.
verbose = true

# directory where  to store downloaded  ads. kleingebaeck will  try to
# create it. must be a quoted string.
outdir = "test"

# template. leave empty to use the default one, which is:
# Title: %s\nPrice: %s\nId: %s\nCategory: %s\nCondition: %s\nCreated: %s\nBody:\n\n%s\n
# take care to include exactly 7 times '%s'!
template = ""
