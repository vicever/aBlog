package theme

import (
	"errors"
	"path/filepath"
	"sync"
)

var (
	tm *themeManager
	l  sync.RWMutex
)

func init() {
	tm = newThemeManager()
}

// themes manager
type themeManager struct {
	current string
	themes  map[string]string
}

func newThemeManager() *themeManager {
	return &themeManager{
		current: "default",
		themes: map[string]string{
			"default": "default",
		},
	}
}

// set theme to global, todo : download and unzip theme, not only add to map
func SetTheme(name string, dir string) {
	if dir == "" {
		delete(tm.themes, name)
		return
	}
	tm.themes[name] = dir
}

// set current theme name,
// if not exist, return error.
func SetCurrent(name string) error {
	l.Lock()
	defer l.Unlock()
	if tm.themes[name] == "" {
		return errors.New("non-exist theme : " + name)
	}
	tm.current = name
	return nil
}

// get current theme name
func GetCurrent(name string) string {
	return tm.current
}

func File(file string) string {
	return filepath.Join(tm.themes[tm.current], file)
}
