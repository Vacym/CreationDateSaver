package internal

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djherbis/times"
	"github.com/fsnotify/fsnotify"
)

// WatchFolder monitors the specified folder and handles file system events.
func WatchFolder(folder string, includeSubfolders bool, metadataFile string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if !filepath.IsAbs(metadataFile) {
		metadataFile = filepath.Join(folder, metadataFile)
	}

	metadata, err := LoadMetadata(metadataFile)
	if err != nil {
		return err
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Skip temporary or unwanted files.
				if IsTemporaryFile(event.Name) {
					continue
				}

				// Calculate relative path once for all operations.
				relPath, err := filepath.Rel(folder, event.Name)
				if err != nil {
					log.Printf("Error calculating relative path: %v\n", err)
					continue
				}

				switch event.Op {
				case fsnotify.Create:
					log.Println("File created:", relPath)

					// Check if it's a directory, if so, add it to the watcher
					if isDir(event.Name) {
						err := watcher.Add(event.Name)
						if err != nil {
							log.Printf("Error adding directory to watcher: %v\n", err)
						}
						log.Printf("New folder added to watcher: %s\n", relPath)
					} else {
						t, err := times.Stat(event.Name)
						if err == nil && t.HasBirthTime() {
							// Use birth time if available.
							UpdateCreationTime(metadata, relPath, t.BirthTime())
						} else {
							// Fall back to modification time if birth time is not available.
							fileInfo, err := os.Stat(event.Name)
							if err == nil {
								UpdateCreationTime(metadata, relPath, fileInfo.ModTime())
							}
						}
					}
					SaveMetadata(metadataFile, metadata)

				case fsnotify.Rename, fsnotify.Remove:
					log.Println("File removed or renamed:", relPath, event.Op)
					DeleteMetadata(metadata, relPath)
					SaveMetadata(metadataFile, metadata)

				case fsnotify.Write:
					log.Println("File written:", relPath)

				default:
					log.Println("Unhandled event:", relPath, event.Op)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	// Add the folder to the watcher.
	err = watcher.Add(folder)
	if err != nil {
		return err
	}

	// Add subfolders if the option is enabled.
	if includeSubfolders {
		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				err := watcher.Add(path)
				if err != nil {
					log.Printf("Error adding directory to watcher: %v\n", err)
				}
			}
			return nil
		})
	}

	<-done
	return nil
}

// Helper function to check if a path is a directory
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
