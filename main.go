package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config := flag.String("config", "config.toml", "path to the config.toml file")
	verbose := flag.Bool("v", false, "enables verbose logging, like showing transactions that didn't match anything")
	interactive := flag.Bool("i", false, "enables interactive mode")
	remove_duplicates := flag.Bool("d", false, "remove duplicates (100.54 | -100.54)")

	flag.Parse()

	conf, err := ParseConfig(*config)
	if err != nil {
		panic(err)
	}

	var matched, unmatched []Transaction
	for _, f := range conf.Files {
		lMatched, lUnmatched := ReadCsv(&conf.ExactRuleMap,
			&conf.PartialRuleMap, &f)
		matched = append(matched, lMatched...)
		unmatched = append(unmatched, lUnmatched...)
	}

	all := append(unmatched, matched...)
	sort.Slice(all, func(i, j int) bool {
		return all[i].Date.Before(all[j].Date)
	})

	if *remove_duplicates {
		discovered := make(map[time.Time]float32)
		for x, i := range all {
			elem, ok := discovered[i.Date]
			if ok && elem == i.Price*-1 {
				if x == len(all) {
					all = all[:x]
				} else {
					all = append(all[:x], all[x+1:]...)
				}
			} else {
				discovered[i.Date] = i.Price
			}
		}
	}

	if *verbose {
		for _, u := range all {
			if u.Destination == "" {
				u.Destination = "UNKNOWN"
			}
			fmt.Println(u)
		}
	}

	if *interactive {
		p := tea.NewProgram(NewState(unmatched))
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
