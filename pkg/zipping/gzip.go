package zipping

import (
	gzipLib "compress/gzip"
	"net/http"
)

type gzip struct {
	rw         http.ResponseWriter
	gzipWriter *gzipLib.Writer
}

func (gzip *gzip) Header() http.Header {
	return gzip.rw.Header()
}

func (gzip *gzip) Write(d []byte) (int, error) {
	return gzip.gzipWriter.Write(d)
}

func (gzip *gzip) WriteHeader(statusCode int) {
	gzip.rw.WriteHeader(statusCode)
}

func (gzip *gzip) Flush() {
	_ = gzip.gzipWriter.Flush()
	_ = gzip.gzipWriter.Close()

}

func NewGzip(rw http.ResponseWriter) *gzip {
	return &gzip{
		gzipWriter: gzipLib.NewWriter(rw),
		rw:         rw,
	}
}
