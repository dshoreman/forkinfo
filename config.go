package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

const configFile = "config.json"
const configPath = ".config/forkinfo"

var config Config

type Config struct {
    AccessToken string `json:"access_token"`
}

func configFullPath() string {
    return strings.Join([] string {
        os.Getenv("HOME"),
        configPath,
        configFile,
    }, "/")
}

func loadConfig() {
    if data, err := ioutil.ReadFile(configFullPath()); err == nil {
        json.Unmarshal(data, &config)
    } else if !os.IsNotExist(err) {
        abortOnError(err)
    }
}

func writeConfig() {
    fmt.Println("Saving config to ", configFullPath(), "...")
    configString, _ := json.MarshalIndent(config, "", "  ")

    os.MkdirAll(strings.Join([] string {os.Getenv("HOME"), configPath}, "/"), 0700)

    if err := ioutil.WriteFile(configFullPath(), append(configString, '\n'), 0644); err != nil {
        fmt.Println("Failed saving config")
        fmt.Println(err)
    }
    loadConfig()
}
