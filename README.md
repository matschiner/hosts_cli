# hosts_cli

This is a utility to easily manage your DNS-entries in /etc/hosts. You can list, add, remove, comment or uncomment entries with it.
====

## Getting Started
For compiling the program Go Lang is required. Please note that it was developed and tested in go1.5.1. It can not be guaranteed that Hosts CLI does also work with other versions of Go.

### Compiling it

For Compiling it into the subdirectory $bin/$ run:

    $ go build -o bin/hosts

If you for example want to add it to your default path for binaries you could run:
    $ go build -o /usr/local/bin/hosts
