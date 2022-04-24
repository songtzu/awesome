// +build !appengine

package lablog

import (
	"io"

	"github.com/mattn/go-colorable"
)

func output() io.Writer {
	return colorable.NewColorableStdout()
}
