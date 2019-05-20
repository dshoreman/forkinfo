package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    if len(os.Args[1:]) == 0 {
        abort("Not enough arguments supplied.")
    }

    if len(os.Args[1:]) > 1 {
        abort("Too many arguments supplied.")
    }

    repo := os.Args[1]

    if ! strings.Contains(repo, "/") {
        abort("Argument is not a valid user/repo string.")
    }

    fmt.Println("Repository set to " + repo)
}

func abort(msg string) {
    fmt.Println("[ERROR] " + msg)
    fmt.Println()
    fmt.Println("Usage: " + os.Args[0] + " <user>/<repository>")
    os.Exit(1)
}
