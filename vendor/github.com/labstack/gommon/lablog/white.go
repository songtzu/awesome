// +build appengine

package lablog

import (
	"io"
	"os"
)

func output() io.Writer {
	return os.Stdout
}
