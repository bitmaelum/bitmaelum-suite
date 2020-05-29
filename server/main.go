package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/server/handler"
    "github.com/jaytaph/mailv2/server/middleware"
    "github.com/sirupsen/logrus"
    "github.com/urfave/negroni"
    "log"
    "net/http"
    "os"
)

var configPath string

func main() {
    parseFlags()
    processConfig()
    processLogging()

    // Main router
    mainRouter := mux.NewRouter().StrictSlash(true)

    mainRouter.HandleFunc("/info", handler.Info).Methods("GET")

    mainRouter.HandleFunc("/account", handler.CreateAccount).Methods("POST")
    mainRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}", handler.RetrieveAccount).Methods("GET")
    mainRouter.HandleFunc("/account/{addr:[A-Za-z0-9]{64}}/key", handler.RetrieveKey).Methods("GET")

    mainRouter.HandleFunc("/incoming", handler.PostMessageHeader).Methods("POST")
    mainRouter.HandleFunc("/incoming/{addr:[A-Za-z0-9]{64}}", handler.PostMessageBody).Methods("POST")


    middlewareRouter := negroni.New()
    middlewareRouter.Use(&middleware.Tracer{})
    middlewareRouter.Use(&middleware.Logger{})
    //middlewareRouter.Use(&middleware.BasicAuth{})
    middlewareRouter.UseHandler(mainRouter)


    host := fmt.Sprintf("%s:%d", config.Server.Server.Host, config.Server.Server.Port)
    logrus.Tracef("listenAndServeTLS on '%s'", host)
    err := http.ListenAndServeTLS(host, config.Server.TLS.CertFile, config.Server.TLS.KeyFile, middlewareRouter)
    if err != nil {
        log.Fatal("listenAndServe: ", err)
    }
}

func processLogging() {
    logrus.SetFormatter(new(logrus.JSONFormatter))
    logrus.SetFormatter(new(logrus.TextFormatter))

    switch (config.Server.Logging.Level) {
    case "trace":
        logrus.SetLevel(logrus.TraceLevel)
        break;
    case "debug":
        logrus.SetLevel(logrus.DebugLevel)
        break;
    case "info":
        logrus.SetLevel(logrus.InfoLevel)
        break;
    case "warning":
        logrus.SetLevel(logrus.WarnLevel)
        break;
    case "error":
    default:
        logrus.SetLevel(logrus.ErrorLevel)
        config.Server.Logging.Level = "error"
        break;
    }
    logrus.SetOutput(os.Stdout)

    logrus.Tracef("setting loglevel to '%s'", config.Server.Logging.Level)
}

func parseFlags() {
    flag.StringVar(&configPath, "config", "./server-config.yml", "path to config file")
    flag.Parse()
}

func processConfig() {
    err := config.Server.LoadConfig(configPath)
    if err != nil {
        panic(err)
    }
}
