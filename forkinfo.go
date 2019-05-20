package main

import (
    "context"
    "fmt"
    "time"
    "os"
    "strconv"
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

func printRepoStats(repo *github.Repository, format string) {
    var datefmt string
    output := "Watchers: %d\tStars: %d\tForks: %d"

    if format == "short" {
        datefmt = "02/01/2006 15:04:05"
        output += "\tIssues/PRs: %d\tLast push: %s\n\n"
    } else {
        datefmt = string(time.RFC1123)
        output += "\n\nThis repository has %d open issues and PRs\nMost recent push: %s\n\n"
    }

    fmt.Printf(
        output,
        repo.GetSubscribersCount(),
        repo.GetStargazersCount(),
        repo.GetForksCount(),
        repo.GetOpenIssuesCount(),
        repo.PushedAt.Format(datefmt),
    )
}

func rowNum(row, total int) string {
    return fmt.Sprintf("[%*d/%d]", len(strconv.Itoa(total)), row, total)
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
    printRepoStats(repo, "long")

    if repo.GetForksCount() == 0 {
        return
    }

    fmt.Printf("Listing forks of %s...\n\n", *repo.FullName)
    forks := fetchRepositoryForks(username, repository)
    numForks := len(forks)

    for i, fork := range forks {
        fmt.Printf("%s %s\n", rowNum(i+1, numForks), *fork.FullName)
        printRepoStats(fork, "short")
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
