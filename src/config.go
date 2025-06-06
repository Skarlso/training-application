package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
)

type appConfig struct {
	configFilePath       string
	alive                bool
	ready                bool
	rootEnabled					 bool
	rootDelaySeconds     int
	startUpDelaySeconds  int
	tearDownDelaySeconds int
	applicationName      string
	applicationVersion   string
	applicationMessage   string
	color                string
	logToFileOnly        bool
	catImageUrl          string
}

func (appConfig *appConfig) logAppConfig() {
	log.Info("Application Configuration:")
	log.Infof("     configFilePath:            %v", appConfig.configFilePath)
	log.Infof("     ready:                     %v", appConfig.ready)
	log.Infof("     alive:                     %v", appConfig.alive)
	log.Infof("     / enabled:                 %v", appConfig.rootEnabled)
	log.Infof("     / delay seconds:           %d", appConfig.rootDelaySeconds)
	log.Infof("     startup delay seconds:     %d", appConfig.startUpDelaySeconds)
	log.Infof("     teardown delay seconds:    %d", appConfig.tearDownDelaySeconds)
	log.Infof("     Application name:          %s", appConfig.applicationName)
	log.Infof("     Applciation version:       %s", appConfig.applicationVersion)
	log.Infof("     Application message:       %s", appConfig.applicationMessage)
	log.Infof("     color:                     %s", appConfig.color)
	log.Infof("     logToFileOnly:             %v", appConfig.logToFileOnly)
	log.Infof("     catImageUrl:               %s", appConfig.catImageUrl)
}

func newAppConfig(configFilePath string) *appConfig {

	ret := &appConfig{
		configFilePath:       configFilePath,
		alive:                true,
		ready:                false,
		rootEnabled: true,
		rootDelaySeconds:     0,
		startUpDelaySeconds:  0,
		tearDownDelaySeconds: 0,
	}

	return ret
}

func (appConfig *appConfig) initAppConfig(isReady bool) {

	appConfig.alive = true
	appConfig.ready = isReady

	fileConfig, err := properties.LoadFile(appConfig.configFilePath, properties.UTF8)
	if err != nil {
		log.Errorf("Configuration file %s not found: 	%v", appConfig.configFilePath, err)
	}

	appConfig.applicationName = getAppConfigStringValue(fileConfig, "name", "APP_NAME", "not set")
	appConfig.applicationVersion = getAppConfigStringValue(fileConfig, "version", "APP_VERSION", "not set")
	appConfig.applicationMessage = getAppConfigStringValue(fileConfig, "message", "APP_MESSAGE", "not set")
	appConfig.color = getAppConfigStringValue(fileConfig, "color", "APP_COLOR", "not set")
	appConfig.logToFileOnly = getAppConfigBoolValue(fileConfig, "logToFileOnly", "", false)
	appConfig.startUpDelaySeconds = getAppConfigIntValue(fileConfig, "startUpDelaySeconds", "", 0)
	appConfig.tearDownDelaySeconds = getAppConfigIntValue(fileConfig, "tearDownDelaySeconds", "", 0)
	catMode := getAppConfigBoolValue(fileConfig, "catMode", "", false)
	if catMode {
		appConfig.catImageUrl, err = getCat()
		if err != nil {
			log.Error("could not obtain cat image", err)
		}
	}
}

func getAppConfigStringValue(fileConfig *properties.Properties, fileConfigProperty, envVarName, defaultValue string) string {
	if envVarName != "" {
		envVarValue, envVarExists := os.LookupEnv(envVarName)
		if envVarExists {
			return envVarValue
		}
	}
	if fileConfig == nil {
		return defaultValue
	}
	return fileConfig.GetString(fileConfigProperty, defaultValue)
}

func getAppConfigBoolValue(fileConfig *properties.Properties, fileConfigProperty, envVarName string, defaultValue bool) bool {
	if envVarName != "" {
		envVarValue, envVarExists := os.LookupEnv(envVarName)
		if envVarExists {
			value, err := strconv.ParseBool(envVarValue)
			if err != nil {
				log.Errorf("could not convert envirnment variable named '%s' with value '%s' to bool:", envVarName, envVarValue)
				return defaultValue
			}
			return value
		}
	}
	if fileConfig == nil {
		return defaultValue
	}
	fileConfigPropertyValue := fileConfig.GetString(fileConfigProperty, "")
	if fileConfigPropertyValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(fileConfigPropertyValue)
	if err != nil {
		log.Errorf("could not convert file configuration property named '%s' with value '%s' to bool:", fileConfigProperty, fileConfigPropertyValue)
		return defaultValue
	} 
	return value
}

func getAppConfigIntValue(fileConfig *properties.Properties, fileConfigProperty, envVarName string, defaultValue int) int {
	if envVarName != "" {
		envVarValue, envVarExists := os.LookupEnv(envVarName)
		if envVarExists {
			value, err := strconv.Atoi(envVarValue)
			if err != nil {
				log.Errorf("could not convert envirnment variable named '%s' with value '%s' to int:", envVarName, envVarValue)
				return defaultValue
			}
			return value
		}
	}
	if fileConfig == nil {
		return defaultValue
	}
	fileConfigPropertyValue := fileConfig.GetString(fileConfigProperty, "")
	if fileConfigPropertyValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(fileConfigPropertyValue)
	if err != nil {
		log.Errorf("could not convert file configuration property named '%s' with value '%s' to int:", fileConfigProperty, fileConfigPropertyValue)
		return defaultValue
	}
	return value
}

func getCat() (string, error) {

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
	err = json.Unmarshal(bodyBytes, &cats)
	if err != nil {
		return "", err
	}
	if len(cats) == 0 {
		return "", errors.New("no cat found")
	}

	return cats[0].Url, nil
}
