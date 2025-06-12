package main

import (
	"bufio"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type cli struct {
	config *appConfig
}

func newCli(appConfig *appConfig) *cli {
	return &cli{
		appConfig,
	}
}

func createHelpText() string {
	var sb strings.Builder
	sb.WriteString("\nAvailable Commands:\n")
	sb.WriteString("\thelp:                get info about available commands and endpoints\n")
	sb.WriteString("\tinit:                set readiness true, liveness true and delay 0\n")
	sb.WriteString("\tconfig:              print out the current application configuration\n")
	sb.WriteString("\tset ready:           application readiness probe will be successful\n")
	sb.WriteString("\tset unready:         application readiness probe will fail\n")
	sb.WriteString("\tset alive:           application liveness probe will be successful\n")
	sb.WriteString("\tset dead:            application liveness probe will fail\n")
	sb.WriteString("\tleak mem:            leak memory\n")
	sb.WriteString("\tleak cpu:            leak cpu\n")
	sb.WriteString("\trequest <url>:       request a url, eg 'request https://www.kubermatic.com/'\n")
	sb.WriteString("\tdelay / <seconds>:   set delay for the root endpoint ('/') in seconds, eg 'delay / 5'\n")
	sb.WriteString("\table Endpoints:\n")
	sb.WriteString("\t/:                   root endpoint, the output is depending on the application configuration\n")
	sb.WriteString("\t/liveness:           liveness probe\n")
	sb.WriteString("\t/readiness:          readiness probe\n")
	return sb.String()
}

func (cli *cli) handleStdin() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("error on reading from stdin: '%s'", err)
		}
		text = strings.ReplaceAll(text, "\n", "")
		if text != "" {
			err = cli.executeCommand(text)
			if err != nil {
				log.Errorf("error on handling command '%s': %s", text, err)
			}
		}
	}
}

func (cli *cli) executeCommand(command string) error {

	if command == "help" {
		log.Info(createHelpText())
	} else if command == "init" {
		log.Info("Re-initializing the application configuration")
		cli.config.initAppConfig(true)
		cli.config.ready = true
		log.Info(cli.config)
	} else if command == "config" {
		log.Info(cli.config)
	} else if command == "set ready" {
		cli.config.ready = true
		log.Info("Set the application to ready")
	} else if command == "set unready" {
		cli.config.ready = false
		log.Info("Set the application to unready")
	} else if command == "set alive" {
		cli.config.alive = true
		log.Info("Set the application to alive")
	} else if command == "set dead" {
		cli.config.alive = false
		log.Info("Set the application to dead")
	} else if command == "leak mem" {
		log.Info("Leaking Memory")
		leakMem()
	} else if command == "leak cpu" {
		log.Info("Leaking CPU")
		leakCpu()
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
		cli.config.rootDelaySeconds, err = strconv.Atoi(delayString)
		if err != nil {
			return fmt.Errorf("error on converting delay string '%s' to int: %s", delayString, err)
		}
		log.Infof("Set delay for the root endpoint ('/') to '%d' seconds", cli.config.rootDelaySeconds)
	} else if strings.HasPrefix(command, "disable /") {
		cli.config.rootEnabled = false
		log.Info("Disabled the root endpoint ('/')")
	} else if strings.HasPrefix(command, "enable /") {
		cli.config.rootEnabled = true
		log.Info("Enabled the root endpoint ('/')")
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("error on closing: %v", err)
		}
	}()

	log.Infof("StatusCode of response %d", resp.StatusCode)

	if resp.TLS == nil {
		log.Info("Response is not encrypted")
		clientCertHeader := resp.Header.Get("X-Client-Cert")
		if clientCertHeader != "" {
			certData, err := base64.StdEncoding.DecodeString(clientCertHeader)
			if err != nil {
				log.Errorf("error decoding proxied certificate: %s", err)
			}
			cert, err := x509.ParseCertificate(certData)
			if err != nil {
				log.Errorf("error parsing proxied certificate: %s", err)
			}
			log.Info(getCertString("Proxied certificate", cert))
		} else {
			log.Infof("No proxied certificate found")
		}
	} else {
		log.Info("Response is encrypted")
		log.Infof("TLS Version: %d", resp.TLS.Version)
		for i, cert := range resp.TLS.PeerCertificates {
			log.Info(getCertString(fmt.Sprintf("Certificate %d", i+1), cert))
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

func getCertString(header string, cert *x509.Certificate) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s: \n", header))
	sb.WriteString(fmt.Sprintf("\tCertificate Subject: %s\n", cert.Subject.String()))
	sb.WriteString(fmt.Sprintf("\tCertificate Issuer: %s\n", cert.Issuer.String()))
	sb.WriteString(fmt.Sprintf("\tCertificate Serial Number: %s\n", cert.SerialNumber.String()))
	sb.WriteString(fmt.Sprintf("\tCertificate Not Before: %s\n", cert.NotBefore.String()))
	sb.WriteString(fmt.Sprintf("\tCertificate Not After: %s\n", cert.NotAfter.String()))
	sb.WriteString(fmt.Sprintf("\tCertificate DNS Names: %v\n", cert.DNSNames))
	sb.WriteString(fmt.Sprintf("\tCertificate Email Addresses: %v\n", cert.EmailAddresses))
	sb.WriteString(fmt.Sprintf("\tCertificate IP Addresses: %v\n", cert.IPAddresses))
	return sb.String()
}

func leakMem() {
	memLeak := make([]string, 0)
	count := 0
	for {
		if count%1000 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.Infof("Alloc = %v MiB", m.Alloc/1024/1024)
			log.Infof("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
			log.Infof("\tSys = %v MiB", m.Sys/1024/1024)
			log.Infof("\tNumGC = %v\n", m.NumGC)
		}
		time.Sleep(time.Nanosecond)
		count++
		memLeak = append(memLeak, "THIS IS A MEM LEAK") //nolint:staticcheck
	}
}

func leakCpu() {

	var waitGroup sync.WaitGroup
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		numGoroutines := runtime.NumGoroutine()
		fmt.Printf("Current number of goroutines: %d\n", numGoroutines)

		// Spawn a new goroutine every second
		waitGroup.Add(1)
		go cpuIntensiveTask(&waitGroup, numGoroutines+1)
	}
}

func cpuIntensiveTask(waitGroup *sync.WaitGroup, id int) {
	defer waitGroup.Done()
	fmt.Printf("Goroutine %d started\n", id)
	for i := 0; i < 1e9; i++ {
		// Perform some CPU-intensive computation
		_ = i * i
	}
	fmt.Printf("Goroutine %d finished\n", id)
}

// // TODO is this really the smartest way to create a CPU leak?

// writer, err := os.Open(os.DevNull)
// if err != nil {
// 	log.Errorf("error on opening /dev/null: %s", err)
// 	return err
// }
// defer writer.Close()
// n := runtime.NumCPU()
// runtime.GOMAXPROCS(n)

// for i := 0; i < n; i++ {
// 	go func() {
// 		for {
// 			var usage syscall.Rusage
// 			err = syscall.Getrusage(syscall.RUSAGE_SELF, &usage)
// 			if err != nil {
// 				log.Errorf("error on cpu usage: %s", err)
// 			}
// 			log.Infof("User CPU Time: %v\n", usage.Utime)
// 			log.Infof("System CPU Time: %v\n", usage.Stime)
// 			fmt.Fprintf(writer, ".")
// 		}
// 	}()
// }

// // TODO do I need this?
// // time.Sleep(10 * time.Second)
// return nil

// }
