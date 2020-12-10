// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler/mgmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware/auth"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/dispatcher"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Configuration file"`
	Version bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	rand.Seed(time.Now().UnixNano())

	internal.ParseOptions(&opts)
	if opts.Version {
		internal.WriteVersionInfo("BitMaelum Server", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	config.LoadServerConfig(opts.Config)
	internal.SetLogging(config.Server.Logging.Format, config.Server.Logging.Level, config.Server.Logging.LogPath)

	checkAndUpdateRouting()

	logrus.Info("Starting " + internal.VersionString("bm-server"))

	// _ = container.Instance.GetResolveService()

	// setup context so we can easily stop all components of the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait for signals and cancel context
	setupSignals(cancel)

	// Start main loop processors
	logrus.Tracef("Starting main loop")
	go mainLoop(ctx)

	// Start webhook dispatcher if enabled
	if config.Server.Webhooks.Enabled {
		initDispatcher()
	}

	// Start HTTP server
	host := fmt.Sprintf("%s:%d", config.Server.Server.Host, config.Server.Server.Port)
	logrus.Tracef("Starting BitMaelum HTTP service on '%s'", host)
	go runHTTPService(ctx, cancel, host)

	// Clean up process queues
	logrus.Tracef("Waiting until context tells us it's done")
	<-ctx.Done()
	logrus.Tracef("Context is done. Exiting")
}

func initDispatcher() {
	logrus.Tracef("Starting webhooks")
	if config.Server.Webhooks.System == "default" {
		logrus.Tracef("Starting default dispatcher system")
		dispatcher.InitDispatcher(config.Server.Webhooks.Workers)
	}
}

func checkAndUpdateRouting() {
	// Check if we have a routing file present
	err := config.ReadRouting(config.Server.Server.RoutingFile)
	if err != nil {
		fmt.Print(`Routing file is not found. We need one in order to uniquely identify this mail server on the network.

You can generate a new one by running:

    $ bm-config generate-routing-id

`)
		os.Exit(1)
	}

	// Check if route exist on the key resolver, and upload new info if needed
	res := container.Instance.GetResolveService()
	info, err := res.ResolveRouting(config.Routing.RoutingID)
	if err != nil || info.Routing != config.Server.Server.Hostname {
		// Upload routing
		err := res.UploadRoutingInfo(resolver.RoutingInfo{
			Hash:      config.Routing.RoutingID,
			PublicKey: config.Routing.PublicKey,
			Routing:   config.Server.Server.Hostname,
		}, &config.Routing.PrivateKey)
		if err != nil {
			fmt.Print("There is an error while uploading routing information to the key resolver: ", err)
			os.Exit(1)
		}
	}
}

func setupSignals(cancel context.CancelFunc) {
	logrus.Tracef("Starting signal notifications")
	go func() {
		// Capture INT and TERM signals
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

		for {
			s := <-sigChannel
			switch s {
			case syscall.SIGHUP:
				// Reload configuration and such
				logrus.Debug("SIGHUP received")
				internal.Reload()
			default:
				logrus.Infof("Signal %s received. Terminating.", s)
				cancel()
			}
		}
	}()
}

// There are two sets of api keys: keys for regular account use, and keys for admin use. They both use the same
// authenticator, but have other permission lists.
var apikeyPermissionList = map[string][]string{
	"ticket": {"send-mail"},
	"boxes":  {"get-header"},
}

var apikeyAdminPermissionList = map[string][]string{
	"admin_flush":  {"flush"},
	"admin_invite": {"invite"},
	"admin_apikey": {"api-key"},
}

func setupRouter() *mux.Router {
	logger := &middleware.Logger{}
	tracer := &middleware.Tracer{}
	prettyJSON := &middleware.PrettyJSON{}

	mainRouter := mux.NewRouter().StrictSlash(true)

	//
	// Router with public endpoints. These can be called without any form of authentication
	//
	publicRouter := mainRouter.PathPrefix("/").Subrouter()
	publicRouter.Use(logger.Middleware)
	publicRouter.Use(prettyJSON.Middleware)
	publicRouter.Use(tracer.Middleware)

	publicRouter.HandleFunc("/", handler.HomePage).Methods("GET")
	publicRouter.HandleFunc("/account", handler.CreateAccount).Methods("POST")
	publicRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/keys", handler.RetrieveKeys).Methods("GET")
	publicRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/org", handler.RetrieveOrganisation).Methods("GET")

	// Server to server message upload
	publicRouter.HandleFunc("/ticket", handler.GetServerToServerTicket).Methods("POST")
	publicRouter.HandleFunc("/incoming/header", handler.IncomingMessageHeader).Methods("POST")
	publicRouter.HandleFunc("/incoming/catalog", handler.IncomingMessageCatalog).Methods("POST")
	publicRouter.HandleFunc("/incoming/block/{message:[A-Za-z0-9-]+}", handler.IncomingMessageBlock).Methods("POST")
	publicRouter.HandleFunc("/incoming/attachment/{message:[A-Za-z0-9-]+}", handler.IncomingMessageAttachment).Methods("POST")
	publicRouter.HandleFunc("/incoming", handler.CompleteIncoming).Methods("POST")
	publicRouter.HandleFunc("/incoming", handler.DeleteIncoming).Methods("DELETE")

	//
	// Routes that need to be authenticated through JWT alone
	//
	authorizer := &middleware.Authenticate{}
	authorizer.Add(&auth.JwtAuth{})

	authRouter := mainRouter.PathPrefix("/").Subrouter()
	authRouter.Use(logger.Middleware)
	authRouter.Use(prettyJSON.Middleware)
	authRouter.Use(tracer.Middleware)
	authRouter.Use(authorizer.Middleware)

	// Message boxes
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/boxes", handler.CreateBox).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[0-9+}", handler.DeleteBox).Methods("DELETE")
	// Message fetching
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/ticket", handler.GetClientToServerTicket).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[0-9]+}", handler.RetrieveMessagesFromBox).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[0-9]+}/message/{message:[A-Za-z0-9-]+}", handler.GetMessage).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[0-9]+}/message/{message:[A-Za-z0-9-]+}/block/{block:[A-Za-z0-9-]+}", handler.GetMessageBlock).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/box/{box:[0-9]+}/message/{message:[A-Za-z0-9-]+}/attachment/{attachment:[A-Za-z0-9-]+}", handler.GetMessageAttachment).Methods("GET")
	// Api keys
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/apikey", handler.CreateAPIKey).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/apikey", handler.ListAPIKeys).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/apikey/{key}", handler.GetAPIKeyDetails).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/apikey/{key}", handler.DeleteAPIKey).Methods("DELETE")
	// Auth keys
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/authkey", handler.CreateAuthKey).Methods("POST")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/authkey", handler.ListAuthKeys).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/authkey/{key}", handler.GetAuthKeyDetails).Methods("GET")
	authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/authkey/{key}", handler.DeleteAuthKey).Methods("DELETE")

	// Add webhooks if enabled
	if config.Server.Webhooks.Enabled {
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook", handler.CreateWebhook).Methods("POST")
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook", handler.ListWebhooks).Methods("GET")
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook/{id}/enable", handler.EnableWebhook).Methods("POST")
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook/{id}/disable", handler.DisableWebhook).Methods("POST")
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook/{id}", handler.GetWebhookDetails).Methods("GET")
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook/{id}", handler.UpdateWebhook).Methods("POST")
		authRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/webhook/{id}", handler.DeleteWebhook).Methods("DELETE")
	}

	//
	// Routes that need to be authenticated through JWT or API keys
	//
	authorizer = &middleware.Authenticate{}
	authorizer.Add(&auth.JwtAuth{})
	authorizer.Add(&auth.APIKeyAuth{
		PermissionList: apikeyPermissionList,
		AdminKeys:      false,
	})

	auth2Router := mainRouter.PathPrefix("/").Subrouter()
	auth2Router.Use(logger.Middleware)
	auth2Router.Use(prettyJSON.Middleware)
	auth2Router.Use(tracer.Middleware)
	auth2Router.Use(authorizer.Middleware)

	// Authorized sending
	auth2Router.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/boxes", handler.RetrieveBoxes).Methods("GET").Name("boxes")

	// Add management endpoints if enabled
	if config.Server.Management.Enabled {
		authorizer = &middleware.Authenticate{}
		authorizer.Add(&auth.APIKeyAuth{
			PermissionList: apikeyAdminPermissionList,
			AdminKeys:      true,
		})

		mgmtRouter := mainRouter.PathPrefix("/admin").Subrouter()
		mgmtRouter.Use(logger.Middleware)
		mgmtRouter.Use(prettyJSON.Middleware)
		mgmtRouter.Use(tracer.Middleware)
		mgmtRouter.Use(authorizer.Middleware)

		mgmtRouter.HandleFunc("/flush", mgmt.FlushQueues).Methods("POST").Name("admin_flush")
		mgmtRouter.HandleFunc("/invite", mgmt.NewInvite).Methods("POST").Name("admin_invite")
		mgmtRouter.HandleFunc("/apikey", mgmt.NewAPIKey).Methods("POST").Name("admin_apikey")
	}

	return mainRouter
}

func runHTTPService(ctx context.Context, cancel context.CancelFunc, addr string) {
	router := setupRouter()

	// Wrap our router in Apache combined logging if needed
	var h http.Handler = router
	if config.Server.Logging.ApacheLogging {
		h = wrapWithApacheLogging(config.Server.Logging.ApacheLogPath, router)
	}

	// Setup HTTP server
	srv := &http.Server{Addr: addr, Handler: h}

	// Start serving TLS in go routine
	go func() {
		err := srv.ListenAndServeTLS(config.Server.Server.CertFile, config.Server.Server.KeyFile)
		if err != nil {
			logrus.Warn("HTTP server stopped: ", err)
			cancel()
		}
	}()

	logrus.Info("HTTP Server up and running")

	// Wait until the context is done
	<-ctx.Done()

	logrus.Info("Shutting down the HTTP server...")
	_ = srv.Close()
}

func mainLoop(ctx context.Context) {
	retryTicker := time.NewTicker(5 * time.Second)
	stuckTicker := time.NewTicker(60 * time.Second)
	processor.IncomingChannel = make(chan string)

	for {
		select {
		// Process incoming message
		case msgID := <-processor.IncomingChannel:
			logrus.Debugf("Message %s uploaded. Processing", msgID)
			err := message.MoveMessage(message.SectionIncoming, message.SectionProcessing, msgID)
			if err != nil {
				logrus.Tracef("issue while trying to move the message %s", err)
				continue
			}
			go processor.ProcessMessage(msgID)

		// Process tickers
		case <-retryTicker.C:
			go processor.ProcessRetryQueue(false)
		case <-stuckTicker.C:
			go processor.ProcessStuckIncomingMessages()
			go processor.ProcessStuckProcessingMessages()

		// Context is done (signal send)
		case <-ctx.Done():
			logrus.Info("Shutting down the processing queues...")
			return
		}
	}
}

// wrapWithApacheLogging wraps a router with a apache logger handler
func wrapWithApacheLogging(path string, router *mux.Router) http.Handler {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		w = os.Stderr
	}

	return handlers.CombinedLoggingHandler(w, router)
}
