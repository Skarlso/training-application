package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	cli "github.com/cloudnativetrainings/training-application/cli"
	conf "github.com/cloudnativetrainings/training-application/conf"
	server "github.com/cloudnativetrainings/training-application/server"

	log "github.com/sirupsen/logrus"
)

var config *conf.AppConfig
var configFilePath string

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {

	configFilePath = getConfigFilePath()

	log.Info("Initializing the application configuration")
	config = conf.NewAppConfig(configFilePath)
	config.InitAppConfig()
	config.LogAppConfig()

	log.Info("Application is starting up")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		log.Infof("Starting the application took %d seconds", i+1)
	}

	cli := cli.NewCli(config)

	go cli.HandleStdin()
	go handleLifecycle()

	server := server.NewServer(config)

	log.Info("Application started, listenting on port 8080")
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
