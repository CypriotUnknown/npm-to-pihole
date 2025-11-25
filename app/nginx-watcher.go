package app

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func (app *AppType) createNginxProxyHostsDirectoryWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add(app.mginxProxyDir)
	if err != nil {
		log.Fatal(err)
	}

	return watcher
}
