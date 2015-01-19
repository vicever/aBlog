package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	logDEBUG = iota
	logINFO
	logWARNING
	logERROR
	logFATAL
)

var (
	logTimeFormat = "01-02 15:04:05"
	logLevels     = [...]string{"Debug", " Info", " Warn", "Error", "Fatal"}
)

type coreLogger struct {
	Writer      io.Writer
	Prefix      string
	NonColor    bool
	ShowDepth   bool
	CallerDepth int
}

func newCoreLogger(file string, prefix string, enableColor bool, enableDepth bool) *coreLogger {
	var w io.Writer = os.Stderr

	// if set to file, use os.File as io.Writer
	if file != "" {
		w, _ = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	}
	l := &coreLogger{
		Writer:      w,
		Prefix:      prefix,
		NonColor:    !enableColor,
		ShowDepth:   enableDepth,
		CallerDepth: 2,
	}

	// windows is not colorful
	if runtime.GOOS == "windows" {
		l.NonColor = true
	}

	return l
}

func (lg *coreLogger) print(level int, format string, args ...interface{}) {
	var depthInfo string
	if lg.ShowDepth {
		pc, file, line, ok := runtime.Caller(lg.CallerDepth)
		if ok {
			// Get caller function name.
			fn := runtime.FuncForPC(pc)
			var fnName string
			if fn == nil {
				fnName = "?()"
			} else {
				fnName = strings.TrimLeft(filepath.Ext(fn.Name()), ".") + "()"
			}
			depthInfo = fmt.Sprintf("[%s:%d %s] ", filepath.Base(file), line, fnName)
		}
	}
	if lg.NonColor {
		fmt.Fprintf(lg.Writer, "%s %s [%s] %s%s\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
		if level == logFATAL {
			os.Exit(1)
		}
		return
	}

	switch level {
	case logDEBUG:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m \033[34m[%s] %s%s\033[0m\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
	case logINFO:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m \033[32m[%s] %s%s\033[0m\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
	case logWARNING:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m \033[33m[%s] %s%s\033[0m\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
	case logERROR:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m \033[31m[%s] %s%s\033[0m\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
	case logFATAL:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m \033[35m[%s] %s%s\033[0m\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
		os.Exit(1)
	default:
		fmt.Fprintf(lg.Writer, "%s %s [%s] %s%s\n",
			lg.Prefix, time.Now().Format(logTimeFormat), logLevels[level], depthInfo,
			fmt.Sprintf(format, args...))
	}
}

func (lg *coreLogger) Debug(format string, args ...interface{}) {
	lg.print(logDEBUG, format, args...)
}

func (lg *coreLogger) Warn(format string, args ...interface{}) {
	lg.print(logWARNING, format, args...)
}

func (lg *coreLogger) Info(format string, args ...interface{}) {
	lg.print(logINFO, format, args...)
}

func (lg *coreLogger) Error(format string, args ...interface{}) {
	lg.print(logERROR, format, args...)
}

func (lg *coreLogger) Fatal(format string, args ...interface{}) {
	lg.print(logFATAL, format, args...)
}
