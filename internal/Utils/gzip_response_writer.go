package utils

import (
	"compress/gzip"
	"net/http"
)

// gzip wrapped response writer
type gzipResponseWrite struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (gzipW *gzipResponseWrite) Write(b []byte) (int, error) {
	return gzipW.writer.Write(b)
}

func NewGzipWriter(r http.ResponseWriter) (writer *gzipResponseWrite, close func()) {
	wr := gzip.NewWriter(r)
	gzWriter := &gzipResponseWrite{
		ResponseWriter: r,
		writer:         wr,
	}
	return gzWriter, func() { wr.Close() }
}
