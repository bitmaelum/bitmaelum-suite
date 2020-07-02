package main

import (
	"context"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-server/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-server/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Configuration file" default:"./server-config.yml"`
	Version bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	if opts.Version {
		core.WriteVersionInfo("BitMaelum Server", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	core.LoadServerConfig(opts.Config)
	core.SetLogging(config.Server.Logging.Level, config.Server.Logging.LogPath)

	logger := &middleware.Logger{}
	tracer := &middleware.Tracer{}
	jwt := &middleware.JwtToken{}

	mainRouter := mux.NewRouter().StrictSlash(true)

	// Public things router
	publicRouter := mainRouter.PathPrefix("/").Subrouter()
	publicRouter.Use(logger.Middleware)
	publicRouter.Use(tracer.Middleware)
	publicRouter.HandleFunc("/", handler.HomePage).Methods("GET")
	publicRouter.HandleFunc("/account", handler.CreateAccount).Methods("POST")
	//publicRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}", handler.RetrieveAccount).Methods("GET")
	publicRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/keys", handler.RetrieveKeys).Methods("GET")
	publicRouter.HandleFunc("/incoming", handler.PostMessageHeader).Methods("POST")
	publicRouter.HandleFunc("/incoming/{addr:[A-Za-z0-9]{64}}", handler.PostMessageBody).Methods("POST")

	//Routes that need to be authenticated
	authRouter := mainRouter.PathPrefix("/").Subrouter()
	authRouter.Use(jwt.Middleware)
	authRouter.Use(logger.Middleware)
	authRouter.Use(tracer.Middleware)
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/boxes", handler.RetrieveBoxes).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[A-Za-z0-9]+}", handler.RetrieveBox).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[A-Za-z0-9]+}/message/{id:[A-Za-z0-9-]+}/flags", handler.GetFlags).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[A-Za-z0-9]+}/message/{id:[A-Za-z0-9-]+}/flag/{flag}", handler.SetFlag).Methods("PUT")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[A-Za-z0-9]+}/message/{id:[A-Za-z0-9-]+}/flag/{flag}", handler.UnsetFlag).Methods("DELETE")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/send/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/header", handler.UploadMessageHeader).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/send/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/catalog", handler.UploadMessageCatalog).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/send/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/block/{id}", handler.UploadMessageBlock).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/send/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", handler.SendMessage).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/send/{uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", handler.DeleteMessage).Methods("DELETE")

	certFilePath, _ := homedir.Expand(config.Server.TLS.CertFile)
	KeyFilePath, _ := homedir.Expand(config.Server.TLS.KeyFile)

	// Wrap our router in Apache combined logging if needed
	var handler http.Handler = mainRouter
	if config.Server.Logging.ApacheLogging == true {
		handler = wrapWithApacheLogging(config.Server.Logging.ApacheLogPath, mainRouter)
	}

	ctx, cancel := context.WithCancel(context.Background())

	logrus.Tracef("Starting processing queues")
	go processQueues(ctx, cancel)

	host := fmt.Sprintf("%s:%d", config.Server.Server.Host, config.Server.Server.Port)
	logrus.Tracef("Starting BitMaelum HTTP service on '%s'", host)
	go runHttpService(ctx, cancel, host, certFilePath, KeyFilePath, handler)

	// Clean up process queues
	logrus.Tracef("Waiting until context is done")
	<-ctx.Done()
	logrus.Tracef("Context is done. Exiting")
}

func runHttpService(ctx context.Context, cancel context.CancelFunc, addr string, certFilePath string, keyFilePath string, h http.Handler) {
	srv := &http.Server{Addr: addr, Handler: h}

	go func() {
		err := srv.ListenAndServeTLS(certFilePath, keyFilePath)
		if err != nil {
			// Cancel context on error
			cancel()
		}
	}()

	logrus.Info("HTTP Server up and running")

	// Wait until the context is done
	for {
		logrus.Info("Waiting for HTTP server to stop...")
		select {
		case <-ctx.Done():
			logrus.Info("Shutting down the HTTP server...")
			_ = srv.Close()
			return
		}
	}
}

func processQueues(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(5 * time.Second)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case uuid := <- processor.UploadChannel:
			fmt.Println("upload message", uuid)
		case uuid := <- processor.OutgoingChannel:
			fmt.Println("outgoing message", uuid)
		case uuid := <- processor.IncomingChannel:
			fmt.Println("incoming message", uuid)
		case t := <-ticker.C:
			fmt.Println("ticker fired at", t)
		case s := <-sigChannel:
			fmt.Println("signal received", s)
			cancel()
		case <-ctx.Done():
			logrus.Info("Shutting down the processing queues...")
			time.Sleep(1 * time.Second)
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("no activity")
		}
	}
}

func wrapWithApacheLogging(path string, router *mux.Router) http.Handler {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		w = os.Stderr
	}

	return handlers.CombinedLoggingHandler(w, router)
}
