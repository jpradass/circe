package pkg

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"

	"github.com/jpradass/circe/db"
	"github.com/jpradass/circe/fs"
)

type Watcher struct {
	filepaths   []string
	destination string
}

func NewWatcher(filepaths []string, destination string) *Watcher {
	return &Watcher{
		filepaths,
		destination,
	}
}

func (w *Watcher) Init() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error when creating fsnotify watcher. Details: %s\n", err.Error())
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("received event: %v\n", event)

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					fileinfo, err := os.Stat(event.Name)
					if err != nil {
						fmt.Printf("there was an error getting stat of file. Details: %s\n", err.Error())
						return
					}

					if fileinfo.IsDir() {
						// look for file inside dir
						fmt.Println("detected directory. Looking for archives inside it...")
					} else {
						if err := db.AddFile(filepath.Base(event.Name), event.Name); err != nil {
							fmt.Printf("there was an error setting file into db. Details: %s\n", err.Error())
							return
						}
						// We move file to destination
						if err := fs.MoveFile(event.Name, w.destination); err != nil {
							fmt.Printf("there was an error moving the file. Details: %s\n", err.Error())
							return
						}

						if err := db.SetStatusToFile(filepath.Base(event.Name)); err != nil {
							fmt.Printf("there was an error settings `sent` status to filename %s", event.Name)
							return
						}
						// When move is done, we need to warn the other part
					}
					// do something
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				fmt.Printf("there was an error on watcher. Details: %s\n", err.Error())
				// print something
			}
		}
	}()

	fmt.Printf("creating watcher for paths: %v\n", w.filepaths)
	for _, path := range w.filepaths {
		err = watcher.Add(path)
		if err != nil {
			return fmt.Errorf("error when adding path %s to watcher. Details: %s", path, err.Error())
		}
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// thread keeps blocked until a signal is received
	sig := <-sigChan

	fmt.Printf("Received signal: %v\n", sig)
	fmt.Println("Exiting...")

	return nil
}
