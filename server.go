package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	config *AppConfig
}

func NewServer(appConfig *AppConfig) *Server {

	server := &Server{
		config: appConfig,
	}

	http.HandleFunc("/", server.handleRoot)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	http.HandleFunc("/liveness", server.handleLiveness)
	http.HandleFunc("/readiness", server.handleReadiness)

	return server
}

func (s *Server) Run() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("Error on starting the server: '%s'", err)
	}
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to root endpoint ('/')")
	if s.config.rootDelay > 0 {
		log.Infof("Delaying for %d seconds", s.config.rootDelay)
		for i := 0; i < s.config.rootDelay; i++ {
			log.Infof("Delayed Response for %d seconds", i+1)
			time.Sleep(1 * time.Second)
		}
		log.Info("Finished delaying Response")
	}
	fmt.Fprintf(w, "<!DOCTYPE html><htlml>")
	fmt.Fprintf(w, "<body style='background-color:%s;'>", s.config.color)
	if s.config.rootDelay > 0 {
		fmt.Fprintf(w, "(Response was delayed for %d seconds)", s.config.rootDelay)
	}
	fmt.Fprintf(w, "Name: %s<br>", s.config.name)
	fmt.Fprintf(w, "Version: %s<br>", s.config.version)
	fmt.Fprintf(w, "Message: %s<br>", s.config.message)
	fmt.Fprintf(w, "Log only to file: %v<br>", s.config.logToFileOnly)
	fmt.Fprintf(w, "Application Liveness: %t<br>", s.config.alive)
	fmt.Fprintf(w, "Application Readiness: %t<br>", s.config.ready)
	fmt.Fprintf(w, "Delay of root endpoint ('/'): %d<br>", s.config.rootDelay)
	fmt.Fprintf(w, "Process ID of the application: %d<br>", os.Getpid())
	if s.config.catImageUrl != "" {
		fmt.Fprintf(w, "<img src='%s' width='500px'></img>", s.config.catImageUrl)
	}
	fmt.Fprintf(w, "</body></htlml>")
}

func (s *Server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to liveness endpoint ('/liveness')")
	if s.config.alive {
		w.WriteHeader(http.StatusOK)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 503 Service Unavailable")
	}
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to readiness endpoint ('/readiness')")
	if s.config.ready {
		w.WriteHeader(http.StatusOK)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 503 Service Unavailable")
	}
}
