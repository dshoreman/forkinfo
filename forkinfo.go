package main

import (
    "context"
    "fmt"
    "time"
    "os"
    "strings"

    "github.com/google/go-github/github"
)

func fetchRepository(username, repository string) (*github.Repository, error) {
    client := github.NewClient(nil)
    data, _, err := client.Repositories.Get(context.Background(), username, repository)

    return data, err
}

func printRepoStats(repo *github.Repository) {
    fmt.Printf(
        "Watchers: %d\tStargazers: %d\tForks: %d\n\n",
        *repo.SubscribersCount,
        *repo.StargazersCount,
        *repo.ForksCount,
    )
    fmt.Println("Most recent push: " + repo.PushedAt.Format(time.RFC1123))
    fmt.Printf("This repository has %d open issues and PRs\n", *repo.OpenIssuesCount)
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

    fmt.Println("Fetching repository...")
    repo, err := fetchRepository(username, repository)
    if (err != nil) {
        fmt.Printf("[ERROR] %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("\n%s: %s\n - %s\n\n", *repo.Name, *repo.Description, *repo.HTMLURL)
    printRepoStats(repo)
}

func abort(msg string) {
    fmt.Println("[ERROR] " + msg)
    fmt.Println()
    fmt.Println("Usage: " + os.Args[0] + " <user>/<repository>")
    os.Exit(1)
}
