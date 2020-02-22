package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

func setupHandlers() {
	http.HandleFunc("/", sayHello)
}

// StartWebServer starts a webserver on the
func StartWebServer(port string) {

	setupHandlers()

	go func() {
		fmt.Printf("serving on %s\n", port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	log.Printf("Webserver running on PORT %s\n", port)
}
