package appender

import (
	"io"
	"os"
)

type Appender struct {
	getWriteCloser func() (io.WriteCloser, error)
}

var _ io.Writer = (*Appender)(nil)

func (a *Appender) Write(v []byte) (int, error) {
	wc, err := a.getWriteCloser()

	if err != nil {
		return 0, err
	}

	defer wc.Close()

	return wc.Write(v)
}

func ForFile(path string) *Appender {
	return &Appender{
		getWriteCloser: func() (io.WriteCloser, error) {
			return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		},
	}
}
