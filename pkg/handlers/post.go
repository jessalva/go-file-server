package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jessalva/go-file-server/pkg/saving"
	"log"
	"mime/multipart"
	"net/http"
)

type PostHandler struct {
	savingService saving.Service
}

func NewPostHandler(savingService saving.Service) *PostHandler {
	return &PostHandler{savingService: savingService}
}

func (ph *PostHandler) SaveFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		filename := vars["filename"]
		postId := vars["id"]

		log.Print(filename + " " + postId)
		err := ph.savingService.SaveFile(filename, postId, r.Body)

		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte("Successfully Saved"))
		if err != nil {
			http.Error(w, "Saved but some issue!", http.StatusInternalServerError)
			return
		}

	}

}

func (ph *PostHandler) SaveFileMultipart() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		postId, file, fileHeader, err, done := extractPostIdAndFile(w, r)
		if done {
			return
		}

		log.Print(fileHeader.Filename + " " + postId)
		err = ph.savingService.SaveFile(fileHeader.Filename, postId, file)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte("Successfully Saved"))
		if err != nil {
			http.Error(w, "Saved but some issue!", http.StatusInternalServerError)
			return
		}

	}

}

func extractPostIdAndFile(w http.ResponseWriter, r *http.Request) (string, multipart.File, *multipart.FileHeader, error, bool) {
	postId := r.FormValue("postId")
	if postId == "" {

		http.Error(w, "Post ID was invalid", http.StatusBadRequest)
		return "", nil, nil, nil, true
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {

		errMsg := fmt.Sprintf("Couldn't get file from form: %v", err)
		http.Error(w, errMsg, http.StatusBadRequest)
		return "", nil, nil, nil, true

	}
	return postId, file, fileHeader, err, false
}
