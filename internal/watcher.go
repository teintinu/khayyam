package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatchCallbacks struct {
	folder   string
	onExec   func()
	onCancel func(canceled func())
	debounce time.Time
	busy     bool
}

var watcher *fsnotify.Watcher
var watchCallbacks []*WatchCallbacks
var watcherMutex sync.Mutex

func InitializeWatcher(onError func(err error)) {
	w, err := fsnotify.NewWatcher()
	if err == nil {
		onError(err)
	} else {
		watcher = w
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					fmt.Println('watching event ', event)
					for _, cb := range watchCallbacks {
						if event.Name == cb.folder {
							watcherMutex.Lock()
							cb.debounce = time.Now().Add(time.Second)
							watcherMutex.Unlock()
						}
					}
				case err := <-watcher.Errors:
					onError(err)
				}
			}
		}()
	}
	go func() {
		for {
			time.Sleep(time.Second)
			watcherMutex.Lock()
			now := time.Now()
			for _, cb := range watchCallbacks {
				if (!cb.debounce.IsZero()) && cb.debounce.After(now) {
					if cb.busy {
						cb.onCancel(func() {
							watcherMutex.Lock()
							cb.busy = false
							watcherMutex.Unlock()
						})
					} else {
						cb.busy = true
						cb.debounce = time.Time{}
						go func() {
							defer func() {
								watcherMutex.Lock()
								cb.busy = false
								watcherMutex.Unlock()
							}()
							cb.onExec()
						}()
					}
				}
			}
			watcherMutex.Unlock()
		}
	}()
}

func WatchFolder(
	folder string,
	onExec func(),
	onCancel func(canceled func()),
	firstRun bool,
) func() {
	watcherMutex.Lock()
	if !isWatching(folder) {
		watcher.Add(folder)
	}
	cb := &WatchCallbacks{
		folder:   folder,
		onExec:   onExec,
		onCancel: onCancel,
		debounce: time.Time{},
	}
	if firstRun {
		cb.debounce = time.Now().Add(time.Second)
	}
	watchCallbacks = append(watchCallbacks, cb)
	watcherMutex.Unlock()
	return func() {
		watcherMutex.Lock()
		var new []*WatchCallbacks
		index := 0
		for _, i := range watchCallbacks {
			if i != cb {
				new = append(new, watchCallbacks[index])
			}
		}
		watchCallbacks = new
		watcherMutex.Unlock()
	}
}

func isWatching(folder string) bool {
	for _, cb := range watchCallbacks {
		if cb.folder == folder {
			return true
		}
	}
	return false
}
