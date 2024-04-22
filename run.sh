#!/bin/sh
# This script is used to compile tailwind's style sheet, and run the web server.
# Also this is supposed to be ran using templ's CLI, i.e.
# templ generate --watch --cmd="./run.sh"
# idk, that's it.

go generate
go run .
