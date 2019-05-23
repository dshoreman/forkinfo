package main

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "time"
    "os"
    "strconv"
    "strings"

    flag "github.com/ogier/pflag"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
)

const configFile = "config.json"
const configPath = ".config/forkinfo"
const version = "0.1.0"

type Config struct {
    AccessToken string `json:"access_token"`
}

var (
    client *github.Client
    config Config
    ctx = context.Background()
    token string
)

func loadConfig() {
    file := strings.Join([] string {
        os.Getenv("HOME"),
        configPath,
        configFile,
    }, "/")

    if data, err := ioutil.ReadFile(file); err == nil {
        json.Unmarshal(data, &config)
    } else if !os.IsNotExist(err) {
        abortOnError(err)
    }
}

func setupAPI() {
    if token = config.AccessToken; token == "" {
        promptForToken()
    }

    token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    client = github.NewClient(oauth2.NewClient(ctx, token))
}

func promptForToken() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("The Github API limits Unauthenticated access to 60 requests per")
    fmt.Println("hour. To raise these limits, create a Personal Access Token at")
    fmt.Println("https://github.com/settings/tokens/new?description=Forkinfo.")
    fmt.Println("Leaves scopes unchecked - Forkinfo requires no special access.")
    fmt.Println()

    for prompt := true; prompt; prompt = token == "" {
        fmt.Println("Paste your API key below:")
        token, _ = reader.ReadString('\n')
        token = strings.Trim(token, " \n\r\t")
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
    var showVersionInfo bool
    flag.BoolVarP(&showVersionInfo, "version", "V", false, "Print version info and quit.")
    flag.Parse()

    if showVersionInfo {
        fmt.Println("Forkinfo " + version)
        os.Exit(0)
    }

    loadConfig()
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
