package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/config"
    "github.com/jaytaph/mailv2/middleware"
    "github.com/jaytaph/mailv2/server"
    "github.com/sirupsen/logrus"
    logger "github.com/sirupsen/logrus"
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
    mainRouter.HandleFunc("/key/{sha256:[A-Za-z0-9]{64}}", server.RetrieveKey).Methods("GET")
    mainRouter.HandleFunc("/key/{sha256:[A-Za-z0-9]{64}}", server.StoreKey).Methods("PUT")
    mainRouter.HandleFunc("/key/{sha256:[A-Za-z0-9]{64}}", server.DeleteKey).Methods("DELETE")

    mainRouter.HandleFunc("/incoming", server.PostMessageHeader).Methods("POST")
    mainRouter.HandleFunc("/incoming/{id:.+}", server.PostMessageBody).Methods("POST")

    // Client router
    // Server router

    middlewareRouter := negroni.New()
    middlewareRouter.Use(&middleware.Logger{})
    //middlewareRouter.Use(&middleware.BasicAuth{})
    middlewareRouter.UseHandler(mainRouter)


    host := fmt.Sprintf("%s:%d", config.Configuration.Server.Host, config.Configuration.Server.Port)
    logger.Tracef("ListenAndServeTLS on '%s'", host)
    err := http.ListenAndServeTLS(host, config.Configuration.TLS.CertFile, config.Configuration.TLS.KeyFile, middlewareRouter)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func processLogging() {
    logger.SetFormatter(new(logrus.JSONFormatter))
    logger.SetFormatter(new(logrus.TextFormatter))
    switch (config.Configuration.Logging.Level) {
    case "trace":
        logger.SetLevel(logrus.TraceLevel)
        break;
    case "debug":
        logger.SetLevel(logrus.DebugLevel)
        break;
    case "info":
        logger.SetLevel(logrus.InfoLevel)
        break;
    case "warning":
        logger.SetLevel(logrus.WarnLevel)
        break;
    case "error":
    default:
        logger.SetLevel(logrus.ErrorLevel)
        config.Configuration.Logging.Level = "error"
        break;
    }
    logger.SetOutput(os.Stdout)

    logger.Tracef("Setting loglevel to '%s'", config.Configuration.Logging.Level)
}

func parseFlags() {
    flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
    flag.Parse()
}

func processConfig() {
    config.Configuration.Logging.Level = "foobar"

    err := config.Configuration.LoadConfig(configPath)
    if err != nil {
        panic(err)
    }
}
