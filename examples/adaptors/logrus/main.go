package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/adaptors"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/stretchr/graceful"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{user}", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page, %s!", mux.Vars(req)["user"])
	})

	mw := interpose.New()

	// Apply the router. By adding it first, all of our other middleware will be
	// executed before the router, allowing us to modify headers before any
	// output has been generated.
	mw.UseHandler(router)

	// Use logrus
	x := negroni.Handler(negronilogrus.NewMiddleware())
	mw.UseHandler(adaptors.FromNegroni(x))

	// Launch and permit graceful shutdown, allowing up to 10 seconds for existing
	// connections to end
	graceful.Run(":3001", 10*time.Second, mw)
}