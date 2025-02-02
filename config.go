package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Rules Rules
	Files []File
}

type Rules struct {
	Exact   []ExactRule
	Partial []PartialRule
}

type ExactRule struct {
	Value   string
	Account string
}

type PartialRule struct {
	Include string
	Account string
}

type File struct {
	Path        string
	Account     string
	DateParse   string
	Date        int
	Description int
	Price       int
}

type ParsedConfig struct {
	ExactRuleMap   map[string]string
	PartialRuleMap map[string]string
	Files          []File
}

func ParseConfig(path string) (ParsedConfig, error) {
	conf := &Config{}
	confFile, err := os.Open(path)
	if err != nil {
		return ParsedConfig{}, err
	}
	defer confFile.Close()
	dec := toml.NewDecoder(confFile)
	err = dec.Decode(conf)
	if err != nil {
		return ParsedConfig{}, err
	}

	descriptionExactMap := make(map[string]string)
	for _, e := range conf.Rules.Exact {
		descriptionExactMap[e.Value] = e.Account
	}

	descriptionPartialMap := make(map[string]string)
	for _, e := range conf.Rules.Partial {
		descriptionPartialMap[e.Include] = e.Account
	}

	return ParsedConfig{
		descriptionExactMap,
		descriptionPartialMap,
		conf.Files,
	}, nil
}
