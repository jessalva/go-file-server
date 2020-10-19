package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jessalva/go-file-server/pkg/saving"
	"github.com/opentracing/opentracing-go"
	"log"
	"mime/multipart"
	"net/http"
)

type PostHandler struct {
	savingService saving.Service
	tracer        opentracing.Tracer
}

func NewPostHandler(savingService saving.Service, tracer opentracing.Tracer) *PostHandler {
	opentracing.SetGlobalTracer(tracer)
	return &PostHandler{savingService: savingService,
		tracer: tracer}
}

func (ph *PostHandler) SaveFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		saveFileSpan, ctx := opentracing.StartSpanFromContextWithTracer(ctx, ph.tracer,"PostHandler::SaveFile")
		defer saveFileSpan.Finish()

		vars := mux.Vars(r)
		filename := vars["filename"]
		postId := vars["id"]

		saveFileSpan.SetBaggageItem("filename", filename)
		saveFileSpan.SetBaggageItem("postId", postId)

		log.Print(filename + " " + postId)
		err := ph.savingService.SaveFile(ctx, filename, postId, r.Body)

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

		ctx := r.Context()
		saveFileSpan, ctx := opentracing.StartSpanFromContextWithTracer(ctx, ph.tracer,"PostHandler::SaveFile")
		defer saveFileSpan.Finish()
		postId, file, fileHeader, err, done := extractPostIdAndFile(w, r)
		if done {
			return
		}

		log.Print(fileHeader.Filename + " " + postId)

		saveFileSpan.SetBaggageItem("filename", fileHeader.Filename)
		saveFileSpan.SetBaggageItem("postId", postId)

		err = ph.savingService.SaveFile(ctx, fileHeader.Filename, postId, file)
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
