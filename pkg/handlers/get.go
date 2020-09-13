package handlers

import (
	"net/http"
	"os"
)

type getHandler struct {}

func NewGetHandler() *getHandler {
	return &getHandler{}
}


func (getHandler *getHandler) GetFile() http.Handler{

	return http.StripPrefix("/images/",http.FileServer( http.Dir(os.Getenv("FILE_SERVER_BASE_PATH")) ))

}
