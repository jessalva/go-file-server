package zipping

import (
	"net/http"
)

type zip interface {
	http.ResponseWriter
}
