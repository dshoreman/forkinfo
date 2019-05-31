package main

import (
    "bufio"
    "context"
    "fmt"
    "net/http"
    "time"
    "os"
    "strconv"
    "strings"

    flag "github.com/ogier/pflag"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
)

const version = "0.1.0"

var (
    authClient *http.Client
    client *github.Client
    ctx = context.Background()
    saveConfig bool
    skipAuth bool
)

func setupAPI() {
    if saveConfig {
        writeConfig()
    }
    if !skipAuth {
        authClient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
            &oauth2.Token{AccessToken: config.AccessToken},
        ))
    }
    client = github.NewClient(authClient)
}

func promptForToken() {
    fmt.Println("The Github API limits Unauthenticated access to 60 requests per")
    fmt.Println("hour. To raise these limits, create a Personal Access Token at")
    fmt.Println("https://github.com/settings/tokens/new?description=Forkinfo.")
    fmt.Println("Leaves scopes unchecked - Forkinfo requires no special access.")
    fmt.Println("To run without authentication, use forkinfo with `--no-token`.")
    fmt.Println()

    reader := bufio.NewReader(os.Stdin)
    for config.AccessToken == "" {
        fmt.Println("Paste your personal access token:")
        token, _ := reader.ReadString('\n')
        config.AccessToken = strings.Trim(token, " \r\n\t")
    }
}

func fetchRepository(username, repository string) (repo *github.Repository) {
    repo, _, err := client.Repositories.Get(ctx, username, repository)
    abortOnError(err)
    return
}

func fetchRepositoryForks(repo *github.Repository) (forks []*github.Repository) {
    opts := github.RepositoryListForksOptions{
        ListOptions: github.ListOptions{PerPage: repo.GetForksCount()},
    }
    forks, _, err := client.Repositories.ListForks(ctx, *repo.Owner.Login, *repo.Name, &opts)
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
    args := strings.Split(flag.Arg(0), "/")
    username, repository := args[0], args[1]

    if config.AccessToken == "" && !skipAuth  {
        promptForToken()
        saveConfig = true
    }
    setupAPI()

    fmt.Println("Fetching repository...")
    repo := fetchRepository(username, repository)

    fmt.Printf("\n%s: %s\n - %s\n\n", *repo.Name, *repo.Description, *repo.HTMLURL)
    printRepoStats(repo, "long")

    if repo.GetForksCount() == 0 {
        return
    }

    fmt.Printf("Listing forks of %s...\n\n", *repo.FullName)
    forks := fetchRepositoryForks(repo)
    numForks := len(forks)

    for i, fork := range forks {
        fmt.Printf("%s %s\n", rowNum(i+1, numForks), *fork.FullName)
        printRepoStats(fork, "short")
    }
}

func init() {
    token := flag.StringP("token", "t", "", "Set the Personal Access Token for API authentication.")
    flag.BoolVarP(&skipAuth, "no-token", "T", false, "Use the Github API without authentication.")
    showVersionInfo := flag.BoolP("version", "V", false, "Print version info and quit.")
    flag.Parse()
    loadConfig()

    if *showVersionInfo {
        fmt.Println("Forkinfo " + version)
        os.Exit(0)
    }
    validate(*token)
    if *token != "" {
        config.AccessToken = *token
    } else if skipAuth {
        fmt.Println("Running without authentication.")
    }
}

func validate(token string) {
    if skipAuth && token != "" {
        abort("Cannot skip authentication while also passing an access token.")
    }
    if len(flag.Args()) == 0 {
        abort("Not enough arguments supplied.")
    }
    if len(flag.Args()) > 1 {
        abort("Too many arguments supplied.")
    }
    if ! strings.Contains(flag.Arg(0), "/") {
        abort("Argument is not a valid user/repo string.")
    }
}

func abort(msg string) {
    fmt.Println("[ERROR] " + msg)
    fmt.Println()
    fmt.Println("Usage: " + os.Args[0] + " <user>/<repository>")
    os.Exit(1)
}

func abortOnError(err error) {
    if err != nil {
        fmt.Printf("[ERROR] %v\n", err)
        os.Exit(1)
    }
}
