package pkg

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
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
				fmt.Printf("received event: %v", event)

				fileinfo, err := os.Stat(event.Name)
				if err != nil {
					fmt.Printf("there was an error getting stat of file. Details: %s\n", err.Error())
					return
				}

				if event.Has(fsnotify.Write) {
					if fileinfo.IsDir() {
						// look for file inside dir
					} else {
						// We move file to destination
						err := os.Rename(event.Name, fmt.Sprintf("%s/%s", w.destination, filepath.Base(event.Name)))
						if err != nil {
							fmt.Printf("there was an error moving file %s to destination %s. Details: %s\n", filepath.Base(event.Name), w.destination, err.Error())
							return
						}

						// When move is done, we need to warn the other part
					}
					// do something
				}

				if event.Has(fsnotify.Create) {
					// do something else
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
