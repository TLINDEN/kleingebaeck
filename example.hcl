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
# "Title: {{.Title}}\nPrice: {{.Price}}\nId: {{.Id}}\nCategory: {{.Category}}\nCondition: {{.Condition}}\nCreated: {{.Created}}\n\n{{.Text}}\n"
template = ""
