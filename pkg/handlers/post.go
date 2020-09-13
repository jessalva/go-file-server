package handlers

import (
	"github.com/gorilla/mux"
	"gitlab.com/jessal.va/go-file-server/pkg/saving"
	"log"
	"net/http"
)

type PostHandler struct{

	savingService saving.Service

}

func NewPostHandler(savingService saving.Service) *PostHandler {
	return &PostHandler{savingService: savingService}
}

func (ph *PostHandler) SaveFile() http.HandlerFunc {

	return func( w http.ResponseWriter, r *http.Request ) {

		vars := mux.Vars( r )
		filename := vars["filename"]
		postId := vars["id"]

		log.Print( filename + " " + postId )
		err := ph.savingService.SaveFile( filename, postId, r.Body )

		if err != nil {

			http.Error( w, err.Error(), http.StatusBadRequest )
			return

		}

		w.Header().Set("Content-Type","application/json")
		_, err = w.Write([]byte("Successfully Saved"))
		if err != nil {
			http.Error( w, "Saved but some issue!", http.StatusInternalServerError )
			return
		}

	}

}