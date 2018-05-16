package auth

import (
	"bufio"
	"context"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// Configuration contains configuration of authorization related information
type Configuration struct {
	TokenFilePath string `yaml:"tokenFile"`
	sync.Mutex
}

func (c *Configuration) getTokenFilePath() string {
	c.Lock()
	defer c.Unlock()
	return c.TokenFilePath
}

func (c *Configuration) setTokenFilePath(tokenFilePath string) {
	c.Lock()
	defer c.Unlock()
	c.TokenFilePath = tokenFilePath
}

var config = &Configuration{}

// SetTokenFilePath sets the auth token file path in the config
func SetTokenFilePath(tokenFilePath string) {
	config.setTokenFilePath(tokenFilePath)
}

type actionList struct {
	actions []func(event fsnotify.Event, logger *zap.Logger)
	sync.Mutex
}

var a = &actionList{}

// AddTokenUpdateAction adds an action to be executed upon changes to the auth file
func AddTokenUpdateAction(action func(event fsnotify.Event, logger *zap.Logger)) {
	a.Lock()
	a.actions = append(a.actions, action)
	a.Unlock()
}

// WatchTokenFile watches an auth token file on disk for changes
func WatchTokenFile(ctx context.Context, logger *zap.Logger) {
	if config.getTokenFilePath() == "" {
		return
	}

	logger.Info("watching updates on " + config.getTokenFilePath())
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	defer watcher.Close()

	err = watcher.Add(config.getTokenFilePath())
	if err != nil {
		if logger != nil {
			logger.Fatal(err.Error())
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("stopping to watch on updates to " + config.getTokenFilePath())
			return
		case event := <-watcher.Events:
			a.Lock()
			for _, action := range a.actions {
				action(event, logger)
			}
			a.Unlock()
		case err := <-watcher.Errors:
			logger.Error(err.Error())
			return
		}
	}
}

// GetToken reads an auth token from file on disk
func GetToken() string {
	return readTokenFromFile(config.getTokenFilePath())
}

func readTokenFromFile(tokenFile string) string {
	file, err := os.Open(tokenFile)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
