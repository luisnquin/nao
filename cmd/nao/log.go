package main

import (
	"io"
	"os"
	"path"

	"github.com/luisnquin/nao/v3/internal"
)

type LogsWriter interface {
	io.StringWriter
	io.Writer
}

func getAppLogsWriter() LogsWriter {
	fp := path.Join(os.TempDir(), "nao.log")
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND

	lf, err := os.OpenFile(fp, flags, internal.PermReadWrite)
	if err != nil {
		panic(err)
	}

	return lf
}
