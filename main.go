package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jpradass/circe/db"
	"github.com/jpradass/circe/pkg"
)

type fp []string

func (f *fp) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *fp) Set(value string) error {
	*f = strings.Split(value, ",")
	return nil
}

func main() {
	var filepaths fp
	destination := flag.String("destination", "", "Destination to place archives")
	flag.Var(&filepaths, "filepaths", "Comma-separated filepaths to watch")

	flag.Parse()

	if len(filepaths) == 0 {
		fmt.Printf("! no filepaths configured. Exiting...\n")
		os.Exit(1)
	}

	// initializes db connection
	if err := db.Init(); err != nil {
		fmt.Printf("error connecting to db: %s", err.Error())
		os.Exit(1)
	}

	watcher := pkg.NewWatcher(filepaths, *destination)
	if err := watcher.Init(); err != nil {
		fmt.Printf("system error: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
