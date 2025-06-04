package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	configFilePath string
	alive          bool
	ready          bool
	rootDelay      int
	name           string
	version        string
	message        string
	color          string
	nodeName       string
	containerName  string
	podNamespace   string
	podName        string
	podIP          string
	logToFileOnly  bool
	catImageUrl    string
}

func (appConfig *AppConfig) LogAppConfig() {
	log.Info("Application Configuration:")
	log.Infof("     configFilePath:  %v", appConfig.configFilePath)
	log.Infof("     ready:           %v", appConfig.ready)
	log.Infof("     alive:           %v", appConfig.alive)
	log.Infof("     / delay:         %d", appConfig.rootDelay)
	log.Infof("     name:            %s", appConfig.name)
	log.Infof("     version:         %s", appConfig.version)
	log.Infof("     message:         %s", appConfig.message)
	log.Infof("     color:           %s", appConfig.color)
	log.Infof("     logToFileOnly:   %v", appConfig.logToFileOnly)
	log.Infof("     nodeName:        %s", appConfig.nodeName)
	log.Infof("     containerName:   %s", appConfig.containerName)
	log.Infof("     podNamespace:    %s", appConfig.podNamespace)
	log.Infof("     podName:         %s", appConfig.podName)
	log.Infof("     podIP:           %s", appConfig.podIP)
	log.Infof("     catImageUrl:     %s", appConfig.catImageUrl)
}

func NewAppConfig(configFilePath string) *AppConfig {

	ret := &AppConfig{
		configFilePath: configFilePath,
		alive:          true,
		ready:          false,
		rootDelay:      0,
	}

	return ret
}

func (appConfig *AppConfig) InitAppConfig() {

	appConfig.alive = true
	appConfig.ready = true
	appConfig.rootDelay = 0

	fileConfig, err := properties.LoadFile(appConfig.configFilePath, properties.UTF8)
	if err != nil {
		log.Errorf("Configuration file %s not found: 	%v", appConfig.configFilePath, err)
	}

	appConfig.name = getAppConfigStringValue(fileConfig, "name", "APP_NAME", "not set")
	appConfig.version = getAppConfigStringValue(fileConfig, "version", "APP_VERSION", "not set")
	appConfig.message = getAppConfigStringValue(fileConfig, "message", "APP_MESSAGE", "not set")
	appConfig.color = getAppConfigStringValue(fileConfig, "color", "APP_COLOR", "not set")
	appConfig.logToFileOnly = getAppConfigBoolValue(fileConfig, "logToFileOnly", "APP_LOG_TO_FILE_ONLY", false)
	appConfig.nodeName = getAppConfigStringValue(fileConfig, "nodeName", "NODE_NAME", "")
	appConfig.containerName = getAppConfigStringValue(fileConfig, "containerName", "CONTAINER_NAME", "")
	appConfig.podNamespace = getAppConfigStringValue(fileConfig, "podNamespace", "POD_NAMESPACE", "")
	appConfig.podName = getAppConfigStringValue(fileConfig, "podName", "POD_NAME", "")
	appConfig.podIP = getAppConfigStringValue(fileConfig, "podIP", "POD_IP", "")
	catMode := getAppConfigBoolValue(fileConfig, "catMode", "APP_CAT_MODE", false)
	if catMode {
		appConfig.catImageUrl, err = getCat()
		if err != nil {
			log.Error("could not obtain cat image", err)
		}
	}
}

func getAppConfigStringValue(fileConfig *properties.Properties, fileConfigProperty, envVarName, defaultValue string) string {
	envVarValue, envVarExists := os.LookupEnv(envVarName)
	if envVarExists {
		return envVarValue
	}
	if fileConfig == nil {
		return defaultValue
	}
	return fileConfig.GetString(fileConfigProperty, "")
}

func getAppConfigBoolValue(fileConfig *properties.Properties, fileConfigProperty, envVarName string, defaultValue bool) bool {
	envVarValue, envVarExists := os.LookupEnv(envVarName)
	if envVarExists {
		value, err := strconv.ParseBool(envVarValue)
		if err != nil {
			log.Errorf("could not convert envirnment variable named '%s' with value '%s' to bool:", envVarName, envVarValue)
			return defaultValue
		}
		return value
	}
	if fileConfig == nil {
		return defaultValue
	}
	fileConfigPropertyValue := fileConfig.GetString(fileConfigProperty, "")
	value, err := strconv.ParseBool(fileConfigPropertyValue)
	if err != nil {
		log.Errorf("could not convert file configuration property named '%s' with value '%s' to bool:", fileConfigProperty, fileConfigPropertyValue)
		return defaultValue
	}
	return value
}

func getCat() (string, error) {

	type catStruct struct {
		url string `json:"url"`
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

	return cats[0].url, nil
}
