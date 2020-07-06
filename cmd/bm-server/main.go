package main

import (
	"context"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
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

	config.LoadServerConfig(opts.Config)
	core.SetLogging(config.Server.Logging.Level, config.Server.Logging.LogPath)

	// setup context so we can easily stop all components of the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait for signals and cancel context
	logrus.Tracef("Starting signal notifications")
	go func() {
		// Capture INT and TERM signals
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

		select {
		case s := <-sigChannel:
			logrus.Print("signal received", s)
			cancel()
		}
	}()

	// Start queues
	logrus.Tracef("Starting processing queues")
	go processQueues(ctx)

	// Start HTTP server
	host := fmt.Sprintf("%s:%d", config.Server.Server.Host, config.Server.Server.Port)
	logrus.Tracef("Starting BitMaelum HTTP service on '%s'", host)
	go runHTTPService(ctx, cancel, host)

	// Clean up process queues
	logrus.Tracef("Waiting until context tells us it's done")
	<-ctx.Done()
	logrus.Tracef("Context is done. Exiting")
}

func setupRouter() *mux.Router {
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

	return mainRouter
}

func runHTTPService(ctx context.Context, cancel context.CancelFunc, addr string) {
	router := setupRouter()

	// Fetch TLS certificate and key
	certFilePath, _ := homedir.Expand(config.Server.TLS.CertFile)
	keyFilePath, _ := homedir.Expand(config.Server.TLS.KeyFile)

	// Wrap our router in Apache combined logging if needed
	var handler http.Handler = router
	if config.Server.Logging.ApacheLogging == true {
		handler = wrapWithApacheLogging(config.Server.Logging.ApacheLogPath, router)
	}

	// Setup HTTP server
	srv := &http.Server{Addr: addr, Handler: handler}

	// Start serving TLS in go routine
	go func() {
		err := srv.ListenAndServeTLS(certFilePath, keyFilePath)
		if err != nil {
			// Cancel context on error
			// @TODO: We should not have cancel here I think, but I don't know a better way to do this.
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

func processQueues(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	processor.UploadChannel = make(chan string)
	processor.OutgoingChannel = make(chan string)
	processor.IncomingChannel = make(chan string)

	for {
		select {
		case uuid := <-processor.UploadChannel:
			fmt.Println("upload message", uuid)
		case uuid := <-processor.OutgoingChannel:
			fmt.Println("outgoing message", uuid)
		case uuid := <-processor.IncomingChannel:
			fmt.Println("incoming message", uuid)
		case t := <-ticker.C:
			fmt.Println("ticker fired at", t)
		case <-ctx.Done():
			logrus.Info("Shutting down the processing queues...")
			time.Sleep(1 * time.Second)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func wrapWithApacheLogging(path string, router *mux.Router) http.Handler {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		w = os.Stderr
	}

	return handlers.CombinedLoggingHandler(w, router)
}
