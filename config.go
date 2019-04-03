package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os/user"
)

var (
	usr, _                  = user.Current()
	home                    = usr.HomeDir
	defaultAlertHandlerPort = ":29000"
	defaultAlertHandlerPath = ""
	defaultConfigPath       = home + "/.alert-handler/alert-handler-config.json"
	defaultScriptsPath      = home + "/.alert-handler/scripts/"
	configPath, scriptPath  *string
)

//Config holds the json structure of an alert handler config file.
type Config struct {
	Port          string                  `json:"port"`
	URLPath       string                  `json:"url-path"`
	ScriptPath    string                  `json:"script-directory"`
	AlertHandlers map[string]AlertHandler `json:"alerts"`
}

//AlertHandler holds the data from an alert as defined in the alert handler config file.
type AlertHandler struct {
	Name       string   `json:"name"`
	Summary    string   `json:"summary"`
	Status     string   `json:"status"`
	ScriptName string   `json:"script-name"`
	ScriptType string   `json:"script-type"`
	Args       []string `json:"args"`
}

func parseConfig(filePath string) (Config, error) {
	file, err := ioutil.ReadFile(filePath)
	config := Config{}
	config.URLPath = defaultAlertHandlerPath
	config.Port = defaultAlertHandlerPort
	config.ScriptPath = defaultScriptsPath
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func init() {

	configPath = flag.String("configPath", defaultConfigPath, "path to configuration file")
	scriptPath = flag.String("scriptPath", defaultScriptsPath, "path to scripts directory")

}

//Set config finds the config file and overwrites the default config values.
func setConfig() Config {
	flag.Parse()
	config, err := parseConfig(*configPath)
	config.ScriptPath = *scriptPath
	if err != nil {
		log.Fatal("Config not parsed correctly:\t", err)
	}
	return config
}
