package main

import (
	"github.com/gorilla/mux"
	"gitlab.com/jessal.va/go-file-server/pkg/handlers"
	"gitlab.com/jessal.va/go-file-server/pkg/saving"
	"gitlab.com/jessal.va/go-file-server/pkg/storage"
	"log"
	"net/http"
	"os"
)

func main() {

	localFileStore := storage.NewLocalFileStore( os.Getenv("FILE_SERVER_BASE_PATH"), 0 )
	savingService := saving.NewService( localFileStore )
	postHandler := handlers.NewPostHandler( savingService )
	getHandler := handlers.NewGetHandler()

	myServeMux := mux.NewRouter()


	getSubRouter := myServeMux.Methods( http.MethodGet ).Subrouter()
	getSubRouter.Handle("/images/{id:[a-zA-Z0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpg|jpeg)}",  getHandler.GetFile())

	putSubRouter := myServeMux.Methods( http.MethodPost ).Subrouter()
	putSubRouter.HandleFunc( "/upload/{id:[a-zA-Z0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpg|jpeg)}", postHandler.SaveFile() )

	server := http.Server{

		Addr: ":8080",
		Handler: myServeMux,

	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Something Bad Happened Yo!")
	}




}
