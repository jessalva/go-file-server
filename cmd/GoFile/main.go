package main

import (
	"github.com/gorilla/mux"
	"github.com/jessalva/go-file-server/pkg/handlers"
	"github.com/jessalva/go-file-server/pkg/saving"
	"github.com/jessalva/go-file-server/pkg/storage"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"net/http"
	"os"
)

func main() {

	if err := run(); err != nil {
		log.Fatal("Something Bad Happened Yo!", err)
	}

}

func run() error {

	logger := log.New( os.Stdout, "FileServer:", log.LstdFlags )

	cfg, err := jaegerCfg.FromEnv()
	if err != nil {
		log.Fatal("cannot parse Jaeger env vars", err.Error())
	}
	cfg.ServiceName = "GO_FILE_SERVER"
	cfg.Sampler.Type = jaeger.SamplerTypeConst
	cfg.Sampler.Param = 1

	jaegerLogger := jaegerLog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, _, err := cfg.NewTracer(
		jaegerCfg.Logger(jaegerLogger),
		jaegerCfg.Metrics(jMetricsFactory),
		jaegerCfg.Observer(rpcmetrics.NewObserver(jMetricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Fatal("cannot initialize Jaeger Tracer",err)
	}

	opentracing.SetGlobalTracer(tracer)

	localFileStore := storage.NewLocalFileStore(os.Getenv("FILE_SERVER_BASE_PATH"), 0, tracer, logger)
	savingService := saving.NewService(localFileStore)
	postHandler := handlers.NewPostHandler(savingService, tracer)
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

		Addr:    ":8090",
		Handler: myServeMux,
	}

	err = server.ListenAndServe()
	return err
}
