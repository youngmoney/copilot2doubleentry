package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	configFilename := flag.String("config", "", "config file (json)")
	// TODO: skip non-cleared items
	var firstDay time.Time
	flag.Func("firstDay", "The first day to include", func(s string) error {
		var err error
		firstDay, err = time.Parse(TRANSACTION_DATE_FORMAT, s)
		return err
	})
	var lastDay time.Time
	flag.Func("lastDay", "The last day to include", func(s string) error {
		var err error
		lastDay, err = time.Parse(TRANSACTION_DATE_FORMAT, s)
		return err
	})
	flag.Parse()
	filename := flag.Arg(0)
	fmt.Printf("first: %s\nlast: %s\nfile: %s\n", firstDay, lastDay, filename)
	ReadCopilot(filename)
	ReadConfig(*configFilename)
}
