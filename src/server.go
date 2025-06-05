package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type server struct {
	config *appConfig
}

func newServer(appConfig *appConfig) *server {

	server := &server{
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

func (s *server) run() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("Error on starting the server: '%s'", err)
	}
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to root endpoint ('/')")
	if s.config.rootDelaySeconds > 0 {
		for i := 0; i < s.config.rootDelaySeconds; i++ {
			log.Infof("Delayed Response for %d of %d seconds", i+1, s.config.rootDelaySeconds)
			time.Sleep(1 * time.Second)
		}
		log.Info("Finished delaying Response")
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<!DOCTYPE html><htlml>")
	fmt.Fprintf(w, "<head><title>%s %s</title></head>", s.config.applicationName, s.config.applicationVersion)
	fmt.Fprintf(w, "<body style='background-color:%s;'>", s.config.color)
	fmt.Fprintf(w, "<h1>%s</h1>", s.config.applicationName)
	fmt.Fprint(w, "<h2>Configuration</h2>")
	fmt.Fprintf(w, "Application Version: %s<br>", s.config.applicationVersion)
	fmt.Fprintf(w, "Application Message: %s<br>", s.config.applicationMessage)
	fmt.Fprintf(w, "Application Liveness: %t<br>", s.config.alive)
	fmt.Fprintf(w, "Application Readiness: %t<br>", s.config.ready)
	fmt.Fprintf(w, "Delay seconds of root endpoint ('/'): %d<br>", s.config.rootDelaySeconds)
	fmt.Fprintf(w, "Seconds the application needs to start up: %d<br>", s.config.startUpDelaySeconds)
	fmt.Fprintf(w, "Seconds the application needs to shut down gracefuly: %d<br>", s.config.tearDownDelaySeconds)
	fmt.Fprintf(w, "Only log to file: %v<br>", s.config.logToFileOnly)

	fmt.Fprint(w, "<h2>Tech Details</h2>")
	fmt.Fprintf(w, "Process Id of the application: %d<br>", os.Getpid())
	fmt.Fprintf(w, "User Id the application is using: %d<br>", os.Getuid())
	hostName, _ := os.Hostname()
	fmt.Fprintf(w, "Hostname: %s<br>", hostName)

	if s.config.catImageUrl != "" {
		fmt.Fprint(w, "<h2>The promised cute cat</h2>")
		fmt.Fprintf(w, "<img src='%s' width='500px'></img>", s.config.catImageUrl)
	}

	fmt.Fprint(w, "</body></htlml>")
}

func (s *server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to liveness endpoint ('/liveness')")
	if s.config.alive {
		w.WriteHeader(http.StatusOK)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 503 Service Unavailable")
	}
}

func (s *server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to readiness endpoint ('/readiness')")
	if s.config.ready {
		w.WriteHeader(http.StatusOK)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 503 Service Unavailable")
	}
}
