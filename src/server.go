package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

//go:embed root.html
var rootTmplContent string

type server struct {
	config *appConfig
	mux    *http.ServeMux
	tmpl   *template.Template
}

type TemplateData struct {
	ApplicationPort      int
	ApplicationName      string
	ApplicationVersion   string
	ApplicationMessage   string
	Color                string
	Alive                bool
	Ready                bool
	RootDelaySeconds     int
	StartUpDelaySeconds  int
	TearDownDelaySeconds int
	RequestInfo          *requestInfo
	LogToFileOnly        bool
	ProcessId            int
	UserId               int
	Hostname             string
	CatImageURL          string
}

func newServer(appConfig *appConfig) *server {

	rootTmpl, err := template.New("root").Parse(rootTmplContent)

	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	mux := http.NewServeMux()

	server := &server{
		config: appConfig,
		mux:    mux,
		tmpl:   rootTmpl,
	}

	mux.HandleFunc("/", server.handleRoot)
	mux.HandleFunc("/favicon.ico", server.handleFavicon)
	mux.HandleFunc("/liveness", server.handleLiveness)
	mux.HandleFunc("/readiness", server.handleReadiness)

	return server
}

func (s *server) run() {
	hostName, _ := os.Hostname()
	log.Infof("Application started with PID %d, UID %d on host with name %s; listenting on port %d", os.Getpid(), os.Getuid(), hostName, config.applicationPort)
	err := http.ListenAndServe(":"+strconv.Itoa(s.config.applicationPort), s.mux)
	if err != nil {
		log.Errorf("error on starting the server: '%s'", err)
	}
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to root endpoint ('/')")
	requestInfo := newRequestInfo(r)
	log.Info(requestInfo)

	if !s.config.rootEnabled {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := fmt.Fprint(w, "The root endpoint of the application is disabled")
		if err != nil {
			log.Errorf("error on writing response for root endpoint ('/'): %s", err)
		}
		log.Info("Root endpoint ('/') responded with Status Code 503 Service Unavailable due to root endpoint is disabled")
		return
	}

	if s.config.rootDelaySeconds > 0 {
		for i := 0; i < s.config.rootDelaySeconds; i++ {
			log.Infof("Delayed Response for %d of %d seconds", i+1, s.config.rootDelaySeconds)
			time.Sleep(1 * time.Second)
		}
		log.Info("Finished delaying Response")
	}

	hostname, _ := os.Hostname()

	data := TemplateData{
		ApplicationPort:      s.config.applicationPort,
		ApplicationName:      s.config.applicationName,
		ApplicationVersion:   s.config.applicationVersion,
		ApplicationMessage:   s.config.applicationMessage,
		Color:                s.config.color,
		Alive:                s.config.alive,
		Ready:                s.config.ready,
		RootDelaySeconds:     s.config.rootDelaySeconds,
		StartUpDelaySeconds:  s.config.startUpDelaySeconds,
		TearDownDelaySeconds: s.config.tearDownDelaySeconds,
		LogToFileOnly:        s.config.logToFileOnly,
		ProcessId:            os.Getpid(),
		UserId:               os.Getuid(),
		RequestInfo:          requestInfo,
		Hostname:             hostname,
		CatImageURL:          s.config.catImageUrl,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := s.tmpl.Execute(w, data); err != nil {
		log.Errorf("error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *server) handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to liveness endpoint ('/liveness')")

	if s.config.alive {
		w.WriteHeader(http.StatusOK)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 500 Internal Server Error")
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


ist es in österreich erlaubt irreführende sponsored links auf google zu schalten, welche auf die falschen webpages verweisen