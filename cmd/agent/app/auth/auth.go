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
	tokenFilePath string `yaml:"tokenFile"`
	token         string
	sync.RWMutex
}

func (c *Configuration) getTokenFilePath() string {
	c.RLock()
	defer c.RUnlock()
	return c.tokenFilePath
}

func (c *Configuration) getToken() string {
	c.RLock()
	defer c.RUnlock()
	return c.token
}

func (c *Configuration) refreshToken() {
	c.Lock()
	defer c.Unlock()
	c.token = readTokenFromFile(c.tokenFilePath)
}

func (c *Configuration) init(tokenFilePath string) {
	c.Lock()
	defer c.Unlock()
	c.tokenFilePath = tokenFilePath
	c.token = readTokenFromFile(tokenFilePath)
}

var config = &Configuration{}

// InitConfig initializes the auth configuration
func InitConfig(tokenFilePath string) {
	config.init(tokenFilePath)
}

type actionList struct {
	actions []func(logger *zap.Logger)
	sync.Mutex
}

var a = &actionList{}

// AddTokenUpdateAction adds an action to be executed upon changes to the auth file
func AddTokenUpdateAction(action func(logger *zap.Logger)) {
	a.Lock()
	a.actions = append(a.actions, action)
	a.Unlock()
}

// WatchTokenChanges watches an auth token file on disk for changes
func WatchTokenChanges(ctx context.Context, logger *zap.Logger) {
	if config.getTokenFilePath() == "" {
		return
	}

	logger.Info("watching content changes on " + config.getTokenFilePath())
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
			if event.Op&fsnotify.Write != fsnotify.Write &&
				event.Op&fsnotify.Create != fsnotify.Create {
				return
			}
			oldToken := config.getToken()
			config.refreshToken()
			newToken := config.getToken()
			if newToken != oldToken {
				logger.Info("[auth] token changed")
				a.Lock()
				for _, action := range a.actions {
					action(logger)
				}
				a.Unlock()
			}
		case err := <-watcher.Errors:
			logger.Error(err.Error())
			return
		}
	}
}

// GetToken reads an auth token from file on disk
func GetToken() string {
	return config.getToken()
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
