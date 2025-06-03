package conf

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
	ConfigFilePath string
	Alive          bool
	Ready          bool
	RootDelay      int
	Name           string
	Version        string
	Message        string
	Color          string
	NodeName       string
	ContainerName  string
	PodNamespace   string
	PodName        string
	PodIP          string
	LogToFileOnly  bool
	CatImageUrl    string
}

func (appConfig *AppConfig) LogAppConfig() {
	log.Info("Application Configuration:")
	log.Infof("     configFilePath:  %v", appConfig.ConfigFilePath)
	log.Infof("     ready:           %v", appConfig.Ready)
	log.Infof("     alive:           %v", appConfig.Alive)
	log.Infof("     / delay:         %d", appConfig.RootDelay)
	log.Infof("     name:            %s", appConfig.Name)
	log.Infof("     version:         %s", appConfig.Version)
	log.Infof("     message:         %s", appConfig.Message)
	log.Infof("     color:           %s", appConfig.Color)
	log.Infof("     logToFileOnly:   %v", appConfig.LogToFileOnly)
	log.Infof("     nodeName:        %s", appConfig.NodeName)
	log.Infof("     containerName:   %s", appConfig.ContainerName)
	log.Infof("     podNamespace:    %s", appConfig.PodNamespace)
	log.Infof("     podName:         %s", appConfig.PodName)
	log.Infof("     podIP:           %s", appConfig.PodIP)
	log.Infof("     catImageUrl:     %s", appConfig.CatImageUrl)
}

func NewAppConfig(configFilePath string) *AppConfig {

	ret := &AppConfig{
		ConfigFilePath: configFilePath,
		Alive:          true,
		Ready:          true,
		RootDelay:      0,
	}

	return ret
}

func (appConfig *AppConfig) InitAppConfig() {

	appConfig.Alive = true
	appConfig.Ready = true
	appConfig.RootDelay = 0

	fileConfig, err := properties.LoadFile(appConfig.ConfigFilePath, properties.UTF8)
	if err != nil {
		log.Errorf("Configuration file %s not found: 	%v", appConfig.ConfigFilePath, err)
	}

	appConfig.Name = getAppConfigStringValue(fileConfig, "name", "APP_NAME", "not set")
	appConfig.Version = getAppConfigStringValue(fileConfig, "version", "APP_VERSION", "not set")
	appConfig.Message = getAppConfigStringValue(fileConfig, "message", "APP_MESSAGE", "not set")
	appConfig.Color = getAppConfigStringValue(fileConfig, "color", "APP_COLOR", "not set")
	appConfig.LogToFileOnly = getAppConfigBoolValue(fileConfig, "logToFileOnly", "APP_LOG_TO_FILE_ONLY", false)
	appConfig.NodeName = getAppConfigStringValue(fileConfig, "nodeName", "NODE_NAME", "")
	appConfig.ContainerName = getAppConfigStringValue(fileConfig, "containerName", "CONTAINER_NAME", "")
	appConfig.PodNamespace = getAppConfigStringValue(fileConfig, "podNamespace", "POD_NAMESPACE", "")
	appConfig.PodName = getAppConfigStringValue(fileConfig, "podName", "POD_NAME", "")
	appConfig.PodIP = getAppConfigStringValue(fileConfig, "podIP", "POD_IP", "")
	catMode := getAppConfigBoolValue(fileConfig, "catMode", "APP_CAT_MODE", false)
	if catMode {
		appConfig.CatImageUrl, err = getCat()
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
