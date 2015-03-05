package theme
import (
    "errors"
    "path/filepath"
)

var tm *themeManager

func init() {
    tm = newThemeManager()
}

// themes manager
type themeManager struct {
    current string
    themes map[string]string
}

func newThemeManager() *themeManager {
    return &themeManager{
        current:"default",
        themes:map[string]string{
            "default":"default",
        },
    }
}

// set theme to global
func SetTheme(name string, dir string) {
    if dir == ""{
        delete(tm.themes,name)
        return
    }
    tm.themes[name] = dir
}

// set current theme name,
// if not exist, return error.
func SetCurrent(name string) error{
    if tm.themes[name] == ""{
        return errors.New("non-exist theme : "+name)
    }
    tm.current = name
    return nil
}

// get current theme name
func GetCurrent(name string) string{
    return tm.current
}

type Themer interface {
    ThemeFile(file string)string
}

type Themed struct {}

func (t *Themed) ThemeFile(file string) string{
    return filepath.Join(tm.themes[tm.current],file)
}