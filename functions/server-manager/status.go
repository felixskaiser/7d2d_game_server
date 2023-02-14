package serverManager

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type statusData struct {
	ServerStatus  string
	HasPlayerInfo bool
	PlayerInfo    string
}

// Check if VM is running and how many players are in the game (implicitely checks if game server is running as well)
func getServerStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving '/status'")

	var status statusData

	isRunning := isVMRunning(projectID, zone, instance)
	if !isRunning {
		status.ServerStatus = "not running"
		status.HasPlayerInfo = false
		getServerStatusRespond(w, status)
		return
	}

	status.ServerStatus = "running"

	players, err := getPlayersInGame(telnetHost, telnetPort, telnetPassword)
	if err != nil {
		status.HasPlayerInfo = false
		getServerStatusRespond(w, status)
		return
	}

	status.HasPlayerInfo = true
	status.PlayerInfo = players

	getServerStatusRespond(w, status)
	return
}

func getServerStatusRespond(w http.ResponseWriter, status statusData) {
	layoutTmplPath := filepath.Join(templateDir, "base.tmpl")
	statusTmplPath := filepath.Join(templateDir, "status.tmpl")

	tmpl, err := template.ParseFiles(layoutTmplPath, statusTmplPath)
	if err != nil {
		log.Printf("SERVER_STATUS: Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", status)
	if err != nil {
		log.Printf("SERVER_STATUS: Error serving '/status': %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	return
}

// Check if VM status is 'RUNNING'
func isVMRunning(projectID, zone, instance string) bool {
	log.Println("SERVER_STATUS: Checking VM status")
	vm, err := computeServiceClient.Instances.Get(projectID, zone, instance).Do()
	if err != nil {
		log.Printf("SERVER_STATUS: Error checking VM status: %v", err)
		return false
	}

	log.Printf("SERVER_STATUS: VM status: %v", vm.Status)
	if vm.Status == "RUNNING" {
		return true
	}

	return false
}

// Query how many players are in the game via game's telnet server
func getPlayersInGame(telnetHost, telnetPort, telnetPassword string) (string, error) {
	log.Printf("SERVER_STATUS: Checking players in game")
	msg, err := sendTelnetCmd("lpi", telnetHost, telnetPort, telnetPassword)
	if err != nil {
		log.Printf("SERVER_STATUS: Error connecting to telnet server for checking players in game: %v", err)
		return "", err
	}

	log.Printf("SERVER_STATUS: Players in game: %s", msg)
	return msg, nil
}
