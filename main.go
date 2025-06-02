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

type FileConfig struct {
	Name    string `properties:"name"`
	Version string `properties:"version"`
	Message string `properties:"message"`
	Color   string `properties:"color"`
	CatMode bool   `properties:"catMode"`
}

type appConfig struct {
	Alive         bool
	Ready         bool
	RootDelay     int
	Name          string
	Version       string
	Message       string
	Color         string
	NodeName      string
	ContainerName string
	PodNamespace  string
	PodName       string
	PodIP         string
	CatImageUrl   string
}

func (appConfig *appConfig) logAppConfig() {
	log.Info("Application Configuration:")
	log.Infof("     ready:           %v", appConfig.Ready)
	log.Infof("     alive:           %v", appConfig.Alive)
	log.Infof("     / delay:         %d", appConfig.RootDelay)
	log.Infof("     name:            %s", appConfig.Name)
	log.Infof("     version:         %s", appConfig.Version)
	log.Infof("     message:         %s", appConfig.Message)
	log.Infof("     color:           %s", appConfig.Color)
	log.Infof("     nodeName:        %s", appConfig.NodeName)
	log.Infof("     containerName:   %s", appConfig.ContainerName)
	log.Infof("     podNamespace:    %s", appConfig.PodNamespace)
	log.Infof("     podName:         %s", appConfig.PodName)
	log.Infof("     podIP:           %s", appConfig.PodIP)
	log.Infof("     catImageUrl:     %s", appConfig.CatImageUrl)
}

var config *appConfig

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Application is starting up")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		log.Infof("Starting the application took %d seconds", i+1)
	}

	log.Info("Initializing the application configuration")
	config = newAppConfig()
}

func newAppConfig() *appConfig {

	var err error

	ret := &appConfig{
		Alive:     true,
		Ready:     true,
		RootDelay: 0,
	}

	fileConfig := properties.MustLoadFile("./conf/app.conf", properties.UTF8)
	ret.Name = initAppConfigValue(fileConfig, "name", "APP_NAME", "not set")
	ret.Version = initAppConfigValue(fileConfig, "version", "APP_VERSION", "not set")
	ret.Message = initAppConfigValue(fileConfig, "message", "APP_MESSAGE", "not set")
	ret.Color = initAppConfigValue(fileConfig, "color", "APP_COLOR", "not set")
	ret.NodeName = initAppConfigValue(fileConfig, "nodeName", "NODE_NAME", "")
	ret.ContainerName = initAppConfigValue(fileConfig, "containerName", "CONTAINER_NAME", "")
	ret.PodNamespace = initAppConfigValue(fileConfig, "podNamespace", "POD_NAMESPACE", "")
	ret.PodName = initAppConfigValue(fileConfig, "podName", "POD_NAME", "")
	ret.PodIP = initAppConfigValue(fileConfig, "podIP", "POD_IP", "")

	catMode := fileConfig.GetBool("catMode", false)
	catModeEnvVarVal, catModeEnvVarValExists := os.LookupEnv("APP_CAT_MODE")
	if catModeEnvVarValExists {
		catMode, err = strconv.ParseBool(catModeEnvVarVal)
		if err != nil {
			log.Errorf("could not convert APP_CAT_MODE '%s' to bool: %s", catModeEnvVarVal, err)
			catMode = false
		}
	}
	if catMode {
		ret.CatImageUrl, err = getCat()
		if err != nil {
			log.Error("could not obtain cat image", err)
		}
	}
	ret.logAppConfig()
	return ret
}

func initAppConfigValue(fileConfig *properties.Properties, fileConfigProperty, envVarName, defaultValue string) string {
	ret := fileConfig.GetString(fileConfigProperty, "")
	envVarValue, envVarExists := os.LookupEnv(envVarName)
	if envVarExists {
		ret = envVarValue
	}
	if ret == "" {
		ret = defaultValue
	}
	return ret
}

func getCat() (string, error) {
	// also string vars => cant they be nil?

	type catStruct struct {
		Url string `json:"url"`
	}

	resp, err := http.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error on reading the response body: '%s'", err)
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Infof("Got response from cat api: %s", bodyString)

	var cats []catStruct
	json.Unmarshal(bodyBytes, &cats)
	if len(cats) == 0 {
		return "", err
	}

	return cats[0].Url, nil
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

	log.Info("Application started, listenting on port 8080")
	log.Info("For getting help, type 'help'")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("Error on starting the server: '%s'", err)
	}
}

func logHelp() {
	log.Info("Available Commands:")
	log.Info("     help:                get info about available commands and endpoints")
	log.Info("     init:                set readiness true, liveness true and delay 0")
	log.Info("     config:              print out the current application configuration")
	log.Info("     set ready:           application readiness probe will be successful")
	log.Info("     set unready:         application readiness probe will fail")
	log.Info("     set alive:           application liveness probe will be successful")
	log.Info("     set dead:            application liveness probe will fail")
	log.Info("     leak mem:            leak memory")
	log.Info("     leak cpu:            leak cpu")
	log.Info("     request <url>:       request a url, eg 'request https://www.google.com'")
	log.Info("     delay / <seconds>:   set delay for the root endpoint ('/') in seconds, eg 'delay / 5'")
	log.Info("Available Endpoints:")
	log.Info("     /:                   root endpoint, the output is depending on the application configuration")
	log.Info("     /liveness:           liveness probe")
	log.Info("     /readiness:          readiness probe")
}

func handleStdin() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("Error on reading from stdin: '%s'", err)
		}
		text = strings.Replace(text, "\n", "", -1)
		if text != "" {
			err = executeCommand(text)
			if err != nil {
				log.Errorf("Error on handling command '%s': %s", text, err)
			}
		}
	}
}

func executeCommand(command string) error {

	if command == "help" {
		logHelp()
	} else if command == "init" {
		log.Info("Re-initializing the application configuration")
		config = newAppConfig()
	} else if command == "config" {
		config.logAppConfig()
	} else if command == "set ready" {
		config.Ready = true
		log.Info("Set the application to ready")
	} else if command == "set unready" {
		config.Ready = false
		log.Info("Set the application to unready")
	} else if command == "set alive" {
		config.Alive = true
		log.Info("Set the application to alive")
	} else if command == "set dead" {
		config.Alive = false
		log.Info("Set the application to dead")
	} else if command == "leak mem" {
		log.Info("Leaking Memory")
		leakMem()
	} else if command == "leak cpu" {
		log.Info("Leaking CPU")
		err := leakCpu()
		if err != nil {
			return fmt.Errorf("error on leaking CPU: %s", err)
		}
	} else if strings.HasPrefix(command, "request ") {
		url, _ := strings.CutPrefix(command, "request ")
		log.Infof("Requesting URL '%s'", url)
		err := request(url)
		if err != nil {
			return fmt.Errorf("error on requesting URL '%s': %s", url, err)
		}
	} else if strings.HasPrefix(command, "delay / ") {
		delayString, _ := strings.CutPrefix(command, "delay / ")
		var err error
		config.RootDelay, err = strconv.Atoi(delayString)
		if err != nil {
			return fmt.Errorf("error on converting delay string '%s' to int: %s", delayString, err)
		}
		log.Infof("Set delay for the root endpoint ('/') to '%d' seconds", config.RootDelay)
	} else {
		return fmt.Errorf("unknown command '%s'", command)
	}
	return nil
}

func request(url string) error {
	log.Infof("Request '%s'", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
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
		return err
	}
	bodyString := string(bodyBytes)
	if len(bodyString) >= 100 {
		bodyString = bodyString[:100]
	}
	log.Infof("Response Body: %s", bodyString)
	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to root endpoint ('/')")
	if config.RootDelay > 0 {
		log.Infof("Delaying for %d seconds", config.RootDelay)
		for i := 0; i < config.RootDelay; i++ {
			log.Infof("Delayed Response for %d seconds", i+1)
			time.Sleep(1 * time.Second)
		}
		log.Info("Finished delaying Response")
	}
	fmt.Fprintf(w, "<!DOCTYPE html><htlml>")
	fmt.Fprintf(w, "<body style='background-color:%s;'>", config.Color)
	if config.RootDelay > 0 {
		fmt.Fprintf(w, "(Response was delayed for %d seconds)", config.RootDelay)
	}
	fmt.Fprintf(w, "Name: %s<br>", config.Name)
	fmt.Fprintf(w, "Version: %s<br>", config.Version)
	fmt.Fprintf(w, "Message: %s<br>", config.Message)
	fmt.Fprintf(w, "Application Liveness: %t<br>", config.Alive)
	fmt.Fprintf(w, "Application Readiness: %t<br>", config.Ready)
	fmt.Fprintf(w, "Delay of root endpoint ('/'): %d<br>", config.RootDelay)
	if config.CatImageUrl != "" {
		fmt.Fprintf(w, "<img src='%s' width='500px'></img>", config.CatImageUrl)
	}
	fmt.Fprintf(w, "</body></htlml>")
}

func handleLiveness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to liveness endpoint ('/liveness')")
	if config.Alive {
		w.WriteHeader(http.StatusOK)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info("Liveness endpoint ('/liveness') responded with Status Code 503 Service Unavailable")
	}
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	log.Info("Request to readiness endpoint ('/readiness')")
	if config.Ready {
		w.WriteHeader(http.StatusOK)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 200 OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Info("Readiness endpoint ('/readiness') responded with Status Code 503 Service Unavailable")
	}
}

func handleLifecycle() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGTERM, syscall.SIGKILL)
	exitChanel := make(chan int)
	go handleSigterm(signalChanel, exitChanel)
	exitCode := <-exitChanel
	os.Exit(exitCode)
}

func handleSigterm(signalChanel chan os.Signal, exitChanel chan int) {
	for {
		signal := <-signalChanel
		if signal == syscall.SIGTERM {
			log.Info("Got SIGTERM signal")
			log.Info("Starting Graceful Shutdown, this will take 10 seconds")
			for i := 0; i < 10; i++ {
				time.Sleep(1 * time.Second)
				log.Infof("Graceful shutdown took %d seconds", i+1)
			}
			log.Info("Graceful Shutdown has finished")
			exitChanel <- 0
		} else if signal == syscall.SIGKILL {
			log.Info("Got SIGKILL signal")
		} else {
			log.Errorf("Got unknown signal '%s'", signal)
			exitChanel <- 1
		}
	}
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

func leakCpu() error {

	// TODO is this really the smartest way to create a CPU leak?

	f, err := os.Open(os.DevNull)
	if err != nil {
		log.Errorf("Error on opening /dev/null: %s", err)
		return err
	}
	defer f.Close()
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	for i := 0; i < n; i++ {
		go func() {
			for {
				var usage syscall.Rusage
				syscall.Getrusage(syscall.RUSAGE_SELF, &usage)
				fmt.Printf("User CPU Time: %v\n", usage.Utime)
				fmt.Printf("System CPU Time: %v\n", usage.Stime)
				fmt.Fprintf(f, ".")
			}
		}()
	}

	// TODO do I need this?
	// time.Sleep(10 * time.Second)
	return nil

}
