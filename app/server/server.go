package main

import (
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/middleware"
    "github.com/jaytaph/mailv2/server"
    "github.com/sirupsen/logrus"
    logger "github.com/sirupsen/logrus"
    "github.com/urfave/negroni"
    "log"
    "net/http"
    "os"
)



func main() {
	logger.SetFormatter(new(logrus.JSONFormatter))
	logger.SetFormatter(new(logrus.TextFormatter))
	logger.SetLevel(logrus.TraceLevel)
	logger.SetOutput(os.Stdout)

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

    err := http.ListenAndServeTLS(":2424", "server.crt", "server.key", middlewareRouter)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
