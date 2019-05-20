package main

import (
    "context"
    "fmt"
    "time"
    "os"
    "strings"

    "github.com/google/go-github/github"
)

var (
    client = github.NewClient(nil)
)

func fetchRepository(username, repository string) (repo *github.Repository) {
    repo, _, err := client.Repositories.Get(context.Background(), username, repository)
    abortOnError(err)
    return
}

func fetchRepositoryForks(username, repository string) (forks []*github.Repository) {
    opts := github.RepositoryListForksOptions{}

    forks, _, err := client.Repositories.ListForks(context.Background(), username, repository, &opts)
    abortOnError(err)
    return
}

func printRepoStats(repo *github.Repository) {
    fmt.Printf(
        "Watchers: %d\tStargazers: %d\tForks: %d\n\n",
        repo.GetSubscribersCount(),
        repo.GetStargazersCount(),
        repo.GetForksCount(),
    )
    fmt.Println("Most recent push: " + repo.PushedAt.Format(time.RFC1123))
    fmt.Printf("This repository has %d open issues and PRs\n\n", *repo.OpenIssuesCount)
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
    repo := fetchRepository(username, repository)

    fmt.Printf("\n%s: %s\n - %s\n\n", *repo.Name, *repo.Description, *repo.HTMLURL)
    printRepoStats(repo)

    if repo.GetForksCount() == 0 {
        return
    }

    fmt.Printf("Listing forks of %s...\n", *repo.FullName)
    forks := fetchRepositoryForks(username, repository)

    for _, fork := range forks {
        fmt.Println()
        fmt.Println(*fork.FullName)
        printRepoStats(fork)
        fmt.Println("---")
    }
}

func abort(msg string) {
    fmt.Println("[ERROR] " + msg)
    fmt.Println()
    fmt.Println("Usage: " + os.Args[0] + " <user>/<repository>")
    os.Exit(1)
}

func abortOnError(err error) {
    if (err != nil) {
        fmt.Printf("[ERROR] %v\n", err)
        os.Exit(1)
    }
}
