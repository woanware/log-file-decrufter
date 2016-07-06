package main

import (
    "github.com/BurntSushi/toml"
    "fmt"
)

// Config struct
type Config struct {
    Misc        misc
    Regex     regex
}

type misc struct {
    ProcessorThreads    int `toml:"processor_threads"`
}

type regex struct {
    Regexes     []string
}


// Load a given configuration by path
func LoadConfig(configPath string) (*Config, error) {
    var config Config
    _, err := toml.DecodeFile(configPath, &config)

    if err != nil {
        return nil, fmt.Errorf("Error loading configuration file: %s", err.Error())
    }

    return &config, nil
}