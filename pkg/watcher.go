package pkg

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	filepaths []string
}

func NewWatcher(filepaths []string) *Watcher {
	return &Watcher{
		filepaths,
	}
}

func (w *Watcher) Init() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error when creating fsnotify watcher. Details: %s", err.Error())
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					// do something
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				fmt.Printf("there was an error on watcher. Details: %s", err.Error())
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
