package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var config *AppConfig
var configFilePath string

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {

	configFilePath = getConfigFilePath()

	log.Info("Initializing the application configuration")
	config = NewAppConfig(configFilePath)
	config.InitAppConfig()
	config.LogAppConfig()

	if config.logToFileOnly {
		log.Warn("Switching to log file only mode, subsequent logs will happen in the file 'application.log'")
		file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	log.Info("Application is starting up")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		log.Infof("Starting the application took %d seconds", i+1)
	}
	config.alive = true
	log.Info("Application set to ready")

	cli := NewCli(config)

	go cli.HandleStdin()
	go handleLifecycle()

	server := NewServer(config)

	log.Infof("Application started with PID %d, listenting on port 8080", os.Getpid())
	log.Info("For getting help, type 'help'")

	server.Run()
}

func getConfigFilePath() string {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "configFilePath" {
		return args[1]
	}
	log.Info("Config File Path not set, defaulting to './conf/app.conf'")
	return "./conf/app.conf"
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
			log.Info("Application set to not ready")
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
