package echo

import (
	"io"

	"github.com/labstack/gommon/lablog"
)

type (
	// Logger defines the logging interface.
	Logger interface {
		Output() io.Writer
		SetOutput(w io.Writer)
		Prefix() string
		SetPrefix(p string)
		Level() lablog.Lvl
		SetLevel(v lablog.Lvl)
		SetHeader(h string)
		Print(i ...interface{})
		Printf(format string, args ...interface{})
		Printj(j lablog.JSON)
		Debug(i ...interface{})
		Debugf(format string, args ...interface{})
		Debugj(j lablog.JSON)
		Info(i ...interface{})
		Infof(format string, args ...interface{})
		Infoj(j lablog.JSON)
		Warn(i ...interface{})
		Warnf(format string, args ...interface{})
		Warnj(j lablog.JSON)
		Error(i ...interface{})
		Errorf(format string, args ...interface{})
		Errorj(j lablog.JSON)
		Fatal(i ...interface{})
		Fatalj(j lablog.JSON)
		Fatalf(format string, args ...interface{})
		Panic(i ...interface{})
		Panicj(j lablog.JSON)
		Panicf(format string, args ...interface{})
	}
)
