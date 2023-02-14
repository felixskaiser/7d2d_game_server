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
	log.Printf("Serving server status")

	var status statusData

	isRunning := isVMRunning(projectID, zone, instance)
	if !isRunning {
		status.ServerStatus = "not running"
		status.HasPlayerInfo = false
	} else {
		status.ServerStatus = "running"
		status.HasPlayerInfo = true

		players, err := getPlayersInGame(telnetHost, telnetPort, telnetPassword)
		if err != nil {
			status.PlayerInfo = "Error connecting to game telnet server for checking players in game, please retry in a few seconds"
		}

		status.PlayerInfo = players
	}

	layoutTmplPath := filepath.Join(templateDir, "base.tmpl")
	statusTmplPath := filepath.Join(templateDir, "status.tmpl")

	tmpl, err := template.ParseFiles(layoutTmplPath, statusTmplPath)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	tmpl.ExecuteTemplate(w, "base", status)
	if err != nil {
		log.Printf("Error serving server status: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	return
}

// Check if VM status is 'RUNNING'
func isVMRunning(projectID, zone, instance string) bool {
	log.Println("VM_STATUS: Checking VM status")
	vm, err := computeServiceClient.Instances.Get(projectID, zone, instance).Do()
	if err != nil {
		log.Printf("VM_STATUS: Error checking VM status: %v", err)
		return false
	}

	log.Printf("VM_STATUS: %v", vm.Status)
	if vm.Status == "RUNNING" {
		return true
	}

	return false
}

// Query how many players are in the game via game's telnet server
func getPlayersInGame(telnetHost, telnetPort, telnetPassword string) (string, error) {
	log.Printf("Checking players in game")
	msg, err := sendTelnetCmd("lpi", telnetHost, telnetPort, telnetPassword)
	if err != nil {
		log.Printf("Error connecting to telnet server for checking players in game: %v", err)
		return "", err
	}

	return msg, nil
}
