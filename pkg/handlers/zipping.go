package handlers

import (
	"github.com/jessalva/go-file-server/pkg/zipping"
	"log"
	"net/http"
	"strings"
)

type zipMiddleware struct {
}

func (z *zipMiddleware) Zip(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		log.Print("Entered Middleware")

		if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {

			nextWriter := zipping.NewDeflate(rw)
			nextWriter.Header().Set("Content-Encoding", "deflate")
			next.ServeHTTP(nextWriter, r)
			defer func() {

				err := nextWriter.Flush()
				if err != nil {
					log.Printf("something bad happened when flushing %v", err)
				}

			}()

			return

		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

			nextWriter := zipping.NewGzip(rw)
			defer nextWriter.Flush()

			nextWriter.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(nextWriter, r)

			return

		}

		next.ServeHTTP(rw, r)

		log.Print("Exited Middleware")

	})

}

func NewZipMiddleWare() *zipMiddleware {

	return &zipMiddleware{}

}
