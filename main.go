package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

const rateLimitPerSecond = 5

var (
	startTime   time.Time
	rateLimiter = time.Tick(1000 / rateLimitPerSecond * time.Millisecond)
)

//Notification is the data structure, in json, actually recieved by the alert reciever. It contains at least one alert.
type Notification struct {
	Alerts []Alert `json:"alerts"`
}

//Alert is the structure of the alertmanager alert which is sent to the server.
type Alert struct {
	Status string            `json:"status"`
	Labels map[string]string `json:"labels"`
}

func getScriptExecCommand(alertHandler AlertHandler) (string, error) {
	switch alertHandler.ScriptType {
	case "python2":
		return "python", nil
	case "bash":
		return "bash", nil
	}
	return "", errors.New("No handler information for script of this type: see config file")
}

//handleAlert looks at the summary of the alert.
//If it finds a match in the list of alertHandler summaries it runs the associated scri:pt.
func handleAlert(alert Alert, alertHandlers map[string]AlertHandler, scriptPath string) error {
	var err error
	var output, scriptExecCommand string
	if handler, ok := alertHandlers[alert.Labels["summary"]]; ok {
		scriptExecCommand, err = getScriptExecCommand(handler)
		if err != nil {
			return err
		}
		script := scriptPath + handler.ScriptName
		args := []string{script}
		for _, arg := range handler.Args {
			args = append(args, arg)
		}
		<-rateLimiter
		fmt.Print(scriptExecCommand)
		_, err := exec.Command(scriptExecCommand, args...).CombinedOutput()
		output = fmt.Sprintf("Firing script %v to deal with alert %v", args[0], alert.Labels["summary"])
		log.Println(output)
		if err != nil {
			return err
		}
		return nil
	}
	err = errors.New("No handler found for this alert: " + alert.Labels["summary"])
	log.Print(err)
	return err
}

func alertReader(config Config, w http.ResponseWriter, r *http.Request) ([]Alert, error) {
	startTime = time.Now()
	alerts := []Alert{}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return alerts, errors.New("Wrong HTTP Method: " + r.Method)
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Read Error:\t%v\n", err.Error())
		return alerts, err
	}
	note := Notification{}
	if err = json.Unmarshal(body, &note); err != nil {
		log.Printf("Json Error:\t%v\n", err.Error())
		return alerts, err
	}

	alerts = note.Alerts
	return alerts, nil
}

func alertsHandler(alerts []Alert, config Config) {
	for _, alert := range alerts {
		go func() {
			handleAlert(alert, config.AlertHandlers, config.ScriptPath)
		}()
	}
}

func handler(config Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alerts, err := alertReader(config, w, r)
		if err != nil {
			log.Printf("%v", err.Error())
			return
		}
		alertsHandler(alerts, config)
	})
}

func main() {
	config := setConfig()
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
		Addr:         config.Port,
	}
	http.HandleFunc(config.URLPath, handler(config))
	log.Printf("Listening on port %v at %v", config.Port, config.URLPath)
	log.Fatal(srv.ListenAndServe())
}
