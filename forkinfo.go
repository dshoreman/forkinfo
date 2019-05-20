package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "github.com/google/go-github/github"
)

func fetchRepository(username, repository string) (*github.Repository, error) {
    client := github.NewClient(nil)
    data, _, err := client.Repositories.Get(context.Background(), username, repository)

    return data, err
}

func main() {
    if len(os.Args[1:]) == 0 {
        abort("Not enough arguments supplied.")
    }
    if len(os.Args[1:]) > 1 {
        abort("Too many arguments supplied.")
    }
    if ! strings.Contains(os.Args[1], "/") {
        abort("Argument is not a valid user/repo string.")
    }

    args := strings.Split(os.Args[1], "/")
    username, repository := args[0], args[1]

    fmt.Println("Attempting to fetch repository...")
    repo, err := fetchRepository(username, repository)
    if (err != nil) {
        fmt.Printf("[ERROR] %v\n", err)
        os.Exit(1)
    }

    json, err := json.MarshalIndent(repo, "", "  ")
    fmt.Print(string(json))
}

func abort(msg string) {
    fmt.Println("[ERROR] " + msg)
    fmt.Println()
    fmt.Println("Usage: " + os.Args[0] + " <user>/<repository>")
    os.Exit(1)
}
