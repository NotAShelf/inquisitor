# Inquisitor

Simple, minimal secrets scanner that I've created for my Forgejo repositories.
Slightly more powerful than a bash script abusing grep with regular expressions
and slightly less than a full fledged Rust program.

## Using

Add "private" headers and path regular expressions to a `configuration.json` and
pass `--config=path/to/configuration.json` to the program.
