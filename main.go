package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
)

var alive = true
var ready = true
var delay = 0

type Config struct {
	Name    string `properties:"name"`
	Version string `properties:"version"`
	Message string `properties:"message"`
	Color   string `properties:"color"`
}

var appConfig *properties.Properties

func init() {
	appConfig = properties.MustLoadFile("./conf/app.conf", properties.UTF8)
}

func main() {

	go handleStdin()
	go handleLifecycle()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	http.HandleFunc("/liveness", handleLiveness)
	http.HandleFunc("/readiness", handleReadiness)
	http.HandleFunc("/downward_api", handleDownwardApi)
	http.HandleFunc("/cats", handleCats)

	log.Info("App started")
	log.Info("AVAILABLE COMMANDS:")
	log.Info("set ready: application readiness probe will be successful")
	log.Info("set unready: application readiness probe will fail")
	log.Info("set alive: application liveness probe will be successful")
	log.Info("set dead: application liveness probe will fail")
	log.Info("leak mem: leak memory")
	log.Info("leak cpu: leak cpu")
	log.Info("request <url>: request a url, eg 'request https://www.google.com'")
	log.Info("delay <seconds>: set delay in seconds, eg 'delay 5'")
	log.Info("AVAILABLE ENDPOINTS:")
	log.Info("/: Root Endpoint, the output is depending on the application configuration")
	log.Info("/liveness: liveness probe")
	log.Info("/readiness: readiness probe")
	log.Info("/downward_api: downward api, giving you metainfo, if available")
	log.Info("/cats: get a random cat image from thecatapi.com")
	http.ListenAndServe(":8080", nil)
}

func handleStdin() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text != "" {
			handleCommand(text)
		}
	}
}

func handleCommand(command string) {
	if command == "set ready" {
		log.Info("Set application to ready")
		ready = true
	} else if command == "set unready" {
		log.Info("Set application to unready")
		ready = false
	} else if command == "set alive" {
		log.Info("Set application to alive")
		alive = true
	} else if command == "set dead" {
		log.Info("Set application to dead")
		alive = false
	} else if command == "leak mem" {
		log.Info("Leaking Memory")
		leakMem()
	} else if command == "leak cpu" {
		log.Info("Leaking CPU")
		leakCpu()
	} else if strings.HasPrefix(command, "request ") {
		url, _ := strings.CutPrefix(command, "request ")
		log.Infof("Requesting URL '%s'", url)
		request(url)
	} else if strings.HasPrefix(command, "delay ") {
		delayString, _ := strings.CutPrefix(command, "delay ")
		delay, _ = strconv.Atoi(delayString)
		// TODO error handling
		log.Infof("Set delay to '%d' seconds", delay)
	} else {
		log.Infof("Unknown command '%s'", command)
	}
}

func request(url string) {
	log.Infof("Request '%s'", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error on getting the response: '%s'", err)
	}
	defer resp.Body.Close()

	log.Infof("StatusCode of response %d", resp.StatusCode)

	if resp.TLS == nil {
		log.Info("Response is not encrypted")
	} else {
		log.Info("Response is encrypted")
		log.Infof("TLS Version: %d", resp.TLS.Version)
		for _, cert := range resp.TLS.PeerCertificates {
			log.Infof("Certificate Subject: %s", cert.Subject.String())
			log.Infof("Certificate Issuer: %s", cert.Issuer.String())
			log.Infof("Certificate Serial Number: %s", cert.SerialNumber.String())
			log.Infof("Certificate Not Before: %s", cert.NotBefore.String())
			log.Infof("Certificate Not After: %s", cert.NotAfter.String())
			log.Infof("Certificate DNS Names: %v", cert.DNSNames)
			log.Infof("Certificate Email Addresses: %v", cert.EmailAddresses)
			log.Infof("Certificate IP Addresses: %v", cert.IPAddresses)
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error on reading the response body: '%s'", err)
	}
	bodyString := string(bodyBytes)
	if len(bodyString) >= 100 {
		bodyString = bodyString[:100]
	}
	log.Infof("Response Body: %s", bodyString)
}

func handleCats(w http.ResponseWriter, r *http.Request) {

	type catStruct struct {
		Url string `json:"url"`
	}

	resp, err := http.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		log.Errorf("Error on getting the response: '%s'", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error on reading the response body: '%s'", err)
	}
	bodyString := string(bodyBytes)
	log.Infof("Got response from cat api: %s", bodyString)

	var cats []catStruct
	json.Unmarshal(bodyBytes, &cats)
	cat := cats[0].Url

	fmt.Fprintf(w, "<!DOCTYPE html><htlml><body>")
	fmt.Fprintf(w, "<img src='%s'></img>", cat)
	fmt.Fprintf(w, "</body></htlml>")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to root")
	if delay > 0 {
		log.Infof("Delaying for %d seconds", delay)
		for i := 0; i < delay; i++ {
			log.Infof("Delayed Response for %d seconds", i+1)
			time.Sleep(1 * time.Second)
		}
		log.Info("Finished delaying Response")
	}
	if ready {
		name := appConfig.GetString("name", "App Configuration Property 'name' is not set")
		version := appConfig.GetString("version", "App Configuration Property 'version' is not set")
		message := appConfig.GetString("message", "App Configuration Property 'message' is not set")
		color := appConfig.GetString("color", "App Configuration Property 'color' is not set")
		fmt.Fprintf(w, "<!DOCTYPE html><htlml>")
		fmt.Fprintf(w, "<body style='background-color:%s;'>", color)
		fmt.Fprintf(w, "Name: %s<br>", name)
		fmt.Fprintf(w, "Version: %s<br>", version)
		fmt.Fprintf(w, "Message: %s<br>", message)
		fmt.Fprintf(w, "Application Liveness: %t<br>", alive)
		fmt.Fprintf(w, "Application Readiness: %t<br>", ready)
		fmt.Fprintf(w, "Application Delay: %d<br>", delay)
		fmt.Fprintf(w, "</body></htlml>")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func handleLiveness(w http.ResponseWriter, r *http.Request) {
	if alive {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func handleDownwardApi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<!DOCTYPE html><htlml><body>")
	fmt.Fprintf(w, "MY_NODE_NAME: %s<br>", os.Getenv("MY_NODE_NAME"))
	fmt.Fprintf(w, "MY_POD_NAME: %s<br>", os.Getenv("MY_POD_NAME"))
	fmt.Fprintf(w, "MY_POD_IP: %s<br>", os.Getenv("MY_POD_IP"))
	fmt.Fprintf(w, "</body></htlml>")
}

func handleLifecycle() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGTERM)
	exitChanel := make(chan int)
	go func() {
		for {
			s := <-signalChanel
			if s == syscall.SIGTERM {
				log.Info("Got SIGTERM signal")
				log.Info("Starting Graceful Shutdown")
				for i := 0; i < 10; i++ {
					log.Infof("Graceful shutdown took %d seconds", i)
					time.Sleep(1 * time.Second)
				}
				log.Info("Graceful Shutdown has finished")
				exitChanel <- 0
			} else {
				log.Info("Got unknown signal")
				exitChanel <- 1
			}
		}
	}()
	exitCode := <-exitChanel
	os.Exit(exitCode)
}

func leakMem() {
	memLeak := make([]string, 0)
	count := 0
	for {
		if count%1000 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
			fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
			fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
			fmt.Printf("\tNumGC = %v\n", m.NumGC)
		}
		time.Sleep(time.Nanosecond)
		count++
		memLeak = append(memLeak, "THIS IS A MEM LEAK")
	}
}

func leakCpu() {

	// TODO is this really the smartest way to create a CPU leak?

	f, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	for i := 0; i < n; i++ {
		go func() {
			for {
				fmt.Fprintf(f, ".")
			}
		}()
	}

	time.Sleep(10 * time.Second)

}
