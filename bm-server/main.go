package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/bitmaelum/bitmaelum-server/bm-server/handler"
    "github.com/bitmaelum/bitmaelum-server/bm-server/middleware"
    "github.com/mitchellh/go-homedir"
    "github.com/sirupsen/logrus"
    "log"
    "net/http"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./server-config.yml"`
}

var opts Options


func main() {
    core.ParseOptions(&opts)
    core.LoadServerConfig(opts.Config)
    core.SetLogging(config.Server.Logging.Level)

    logger := &middleware.Logger{}
    tracer := &middleware.Tracer{}
    jwt := &middleware.JwtToken{}

    mainRouter := mux.NewRouter().StrictSlash(true)

    // Public things router
    publicRouter := mainRouter.PathPrefix("/").Subrouter()
    publicRouter.Use(logger.Middleware)
    publicRouter.Use(tracer.Middleware)
    publicRouter.HandleFunc("/info", handler.Info).Methods("GET")
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
    authRouter.HandleFunc("/send", handler.SendMessage).Methods("POST")


    cfp, _ := homedir.Expand(config.Server.TLS.CertFile)
    kfp, _ := homedir.Expand(config.Server.TLS.KeyFile)

    host := fmt.Sprintf("%s:%d", config.Server.Server.Host, config.Server.Server.Port)
    logrus.Tracef("listenAndServeTLS on '%s'", host)
    err := http.ListenAndServeTLS(host, cfp, kfp, mainRouter)
    if err != nil {
        log.Fatal("listenAndServe: ", err)
    }
}
