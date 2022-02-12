package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Nirespire/twitchbot/types"
)

type webServer interface {
	StartWebServer()
	handleConfig()
	setupHandlers()
	sayHello()
}

type ServerConfig struct {
	BotConfig *types.ChatConfig
	BotStats *types.BotStats
	Port string
}

type configRequest struct {
	Name string
	Value string
}

func (config *ServerConfig) sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

// TODO Add logging
func (config *ServerConfig) handleConfig(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		jsonConfig, err := json.Marshal(config)

		if err != nil {
			http.Error(w, "Error parsing Bot Configuration",
					http.StatusInternalServerError)
		}
	
		w.Write([]byte (jsonConfig))

	} else if r.Method == "POST" {

		var req configRequest

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		
		err = json.Unmarshal([]byte(body), &req)

		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}

		switch req.Name {
		case "ProjectDescription":
			config.BotConfig.ProjectDescription = req.Value

			log.Printf("Changed project message to %s", config.BotConfig.ProjectDescription)
		default:
			log.Printf("Invalid config Name: %s", req.Name)
			http.Error(w, "Invalid config Name",
				http.StatusBadRequest)
			return
		}

		jsonConfig, err := json.Marshal(config)

		if err != nil {
			http.Error(w, "Error parsing Bot Configuration",
					http.StatusInternalServerError)
		}
	
		w.Write([]byte (jsonConfig))

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func (server *ServerConfig) setupHandlers() {
	http.HandleFunc("/", server.sayHello)
	http.HandleFunc("/config", server.handleConfig)
}

// StartWebServer starts a webserver on the
func (server *ServerConfig) StartWebServer() {

	server.setupHandlers()

	go func() {
		fmt.Printf("serving on %s\n", server.Port)
		err := http.ListenAndServe(server.Port, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	log.Printf("Webserver running on PORT %s\n", server.Port)
}
