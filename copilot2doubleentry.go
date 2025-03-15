package main

import (
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

func main() {
	configFilename := flag.String("config", "", "config file (json)")
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
	if len(filename) == 0 {
		fmt.Fprintln(os.Stderr, "provide a csv file to convert")
		os.Exit(1)
	}
	transactions := ReadCopilot(filename)
	config := ReadConfig(*configFilename)
	converted := Convert(transactions, config, firstDay, lastDay)
	csvContent, err := gocsv.MarshalString(&converted)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(csvContent)
}
