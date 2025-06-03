package server

import (
	"fmt"
	"net/http"
	"time"

	conf "github.com/cloudnativetrainings/training-application/conf"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config *conf.AppConfig
}

func NewServer(appConfig *conf.AppConfig) *Server {

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

func (s Server) Run() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("Error on starting the server: '%s'", err)
	}

}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to root endpoint ('/')")
	if s.config.RootDelay > 0 {
		log.Infof("Delaying for %d seconds", s.config.RootDelay)
		for i := 0; i < s.config.RootDelay; i++ {
			log.Infof("Delayed Response for %d seconds", i+1)
			time.Sleep(1 * time.Second)
		}
		log.Info("Finished delaying Response")
	}
	fmt.Fprintf(w, "<!DOCTYPE html><htlml>")
	fmt.Fprintf(w, "<body style='background-color:%s;'>", s.config.Color)
	if s.config.RootDelay > 0 {
		fmt.Fprintf(w, "(Response was delayed for %d seconds)", s.config.RootDelay)
	}
	fmt.Fprintf(w, "Name: %s<br>", s.config.Name)
	fmt.Fprintf(w, "Version: %s<br>", s.config.Version)
	fmt.Fprintf(w, "Message: %s<br>", s.config.Message)
	fmt.Fprintf(w, "LogToFileOnly: %v<br>", s.config.LogToFileOnly)
	fmt.Fprintf(w, "Application Liveness: %t<br>", s.config.Alive)
	fmt.Fprintf(w, "Application Readiness: %t<br>", s.config.Ready)
	fmt.Fprintf(w, "Delay of root endpoint ('/'): %d<br>", s.config.RootDelay)
	if s.config.CatImageUrl != "" {
		fmt.Fprintf(w, "<img src='%s' width='500px'></img>", s.config.CatImageUrl)
	}
	fmt.Fprintf(w, "</body></htlml>")
}

func (s *Server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to liveness endpoint ('/liveness')")
	if s.config.Alive {
		w.WriteHeader(http.StatusOK)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 503 Service Unavailable")
	}
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to readiness endpoint ('/readiness')")
	if s.config.Ready {
		w.WriteHeader(http.StatusOK)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 503 Service Unavailable")
	}
}
