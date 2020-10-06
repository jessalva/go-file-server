package zipping

import (
	"compress/flate"
	"net/http"
)

type deflate struct {
	rw            http.ResponseWriter
	deflateWriter *flate.Writer
}

func (d *deflate) Header() http.Header {
	return d.rw.Header()
}

func (d *deflate) Write(bytes []byte) (int, error) {
	return d.deflateWriter.Write(bytes)
}

func (d *deflate) WriteHeader(statusCode int) {
	d.rw.WriteHeader(statusCode)
}

func (d *deflate) Flush() error {
	return d.deflateWriter.Flush()
}

func NewDeflate(rw http.ResponseWriter) *deflate {

	deflateWriter, _ := flate.NewWriter(rw, 1)
	return &deflate{
		deflateWriter: deflateWriter,
		rw:            rw,
	}
}
