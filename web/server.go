package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Nirespire/twitchbot/util"
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
	fmt.Printf("[%s] Webserver running on PORT %s\n", util.TimeStamp(), port)
}
