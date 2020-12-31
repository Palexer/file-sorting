package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	// parse flag for directory
	dir := flag.String("dir", ".", "specify a directory to sort")
	flag.Parse()

	log.Printf("watching for changes in %s\n", *dir)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	done := make(chan bool)

	// look for events
	go func() {
		for {
			select {
			// watch for events
			case event := <-w.Events:
				log.Printf("event: %#v\n", event) // log events
				if event.Op == fsnotify.Create {  // only for new files
					if filepath.Ext(event.Name) != "" && !strings.HasPrefix(filepath.Base(event.Name), ".") { // filter out folders and hidden files
						// get file modification time
						info, _ := os.Stat(event.Name)
						statt := info.Sys().(*syscall.Stat_t)
						t := timespecToTime(statt.Mtim)

						year := strconv.Itoa(t.Year())
						month := strconv.Itoa(int(t.Month()))

						// 0 in folder name for months < 10
						if len(month) <= 1 {
							month = "0" + month
						}

						newPath := *dir + "/" + year + "/" + month + "/" + filepath.Base(event.Name)

						log.Printf("\n\n%s\n\n", newPath)

						// create directories if they don't exist
						// year
						if e, _ := exists(*dir + "/" + year); !e {
							os.Mkdir(*dir+"/"+year, os.ModePerm)
						}
						// month
						if e, _ := exists(*dir + "/" + year + "/" + month); !e {
							os.Mkdir(*dir+"/"+year+"/"+month, os.ModePerm)
						}

						// move file into the new directory
						err := os.Rename(event.Name, newPath)
						if err != nil {
							log.Printf("failed to move file: %v", err)
						}
					}

				}

			// watch for errors
			case err := <-w.Errors:
				log.Printf("error: %v", err)
			}
		}
	}()

	// Watch  *dir for changes.
	if err := w.Add(*dir); err != nil {
		log.Fatal(err)
	}
	<-done
}
