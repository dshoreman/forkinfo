package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args[1:]) != 1 {
        fmt.Println("Not enough or too many arguments supplied.")
        fmt.Println()
        fmt.Println("Usage: " + os.Args[0] + " <user>/<repository>")
        os.Exit(1)
    }

    repo := os.Args[1]

    fmt.Println("Repository set to " + repo)
}
