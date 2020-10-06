package main

import (
	"github.com/gorilla/mux"
	"github.com/jessalva/go-file-server/pkg/handlers"
	"github.com/jessalva/go-file-server/pkg/saving"
	"github.com/jessalva/go-file-server/pkg/storage"
	"log"
	"net/http"
	"os"
)

func main() {

	localFileStore := storage.NewLocalFileStore(os.Getenv("FILE_SERVER_BASE_PATH"), 0)
	savingService := saving.NewService(localFileStore)
	postHandler := handlers.NewPostHandler(savingService)
	getHandler := handlers.NewGetHandler()
	zipMiddleWare := handlers.NewZipMiddleWare()

	myServeMux := mux.NewRouter()

	getSubRouter := myServeMux.Methods(http.MethodGet).Subrouter()
	getSubRouter.Handle("/images/{id:[a-zA-Z0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpg|jpeg|JPG)}", getHandler.GetFile())
	getSubRouter.Use(zipMiddleWare.Zip)

	postSubRouter := myServeMux.Methods(http.MethodPost).Subrouter()
	postSubRouter.HandleFunc("/upload/{id:[a-zA-Z0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpg|jpeg)}", postHandler.SaveFile())
	postSubRouter.HandleFunc("/", postHandler.SaveFileMultipart())

	server := http.Server{

		Addr:    ":8080",
		Handler: myServeMux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Something Bad Happened Yo!")
	}

}
