# forkinfo

Github will show you a list of forks for any repository via its Insights page, but to find out any other
useful information such as when the fork was last updated, how many commits differ etc., you're
forced to open each and every fork in a new tab... Often to find that most haven't changed at all!

Forkinfo aims to fix that by providing basic information about all forks with one simple command.

## Installation
```
go get github.com/dshoreman/forkinfo
```

### Development

To build Forkinfo without installing globally, clone this repo and run `go build` in the project directory.
This will compile the application to `forkinfo` which can be run as `./forkinfo <user>/<repo>`.

Alternatively, you can use `go run forkinfo.go` which will compile into `/tmp` and run it automatically.

## Usage

Forkinfo takes a single argument, and that's the `user/repo` string. For example, to check this repo:

```
$ forkinfo dshoreman/forkinfo
```

When you run Forkinfo you will be prompted to enter a Personal Access Token to avoid hitting rate
limits on the Github API. This is optional, but will be saved to a config file once you paste in your
token and hit enter. To skip authentication, run Forkinfo with the `--no-token` flag. Note that "guest"
API access is limited to 60 requests per hour, which may be too low for some repositories.
