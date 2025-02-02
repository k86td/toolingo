package main

import (
	"flag"
	"fmt"
)

func main() {
	config := flag.String("config", "config.toml", "path to the config.toml file")
	verbose := flag.Bool("v", false, "enables verbose logging, like showing transactions that didn't match anything")

	flag.Parse()

	conf, err := ParseConfig(*config)
	if err != nil {
		panic(err)
	}

	matched, unmatched := ReadCsv(&conf,
		*fileName, *account)

	for _, m := range matched {
		fmt.Println(m)
	}

	if *verbose {
		for _, u := range unmatched {
			u.Destination = "UNKNOWN"
			fmt.Println(u)
		}
	}
}
