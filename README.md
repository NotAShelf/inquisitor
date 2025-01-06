# Inquisitor

Simple, minimal secrets scanner that I've created for my Forgejo repositories.
Slightly more powerful than a bash script abusing grep with regular expressions
and slightly less than a full fledged Rust program.

## Using

There is no configurability, so you will need to edit `main.go` to add your
regular expressions to the list. Build with `go build`, and install it anywhere
in your PATH. Running `inquisitor` will scan the current directory for certain
headers that indicate sensitive content in the file.
