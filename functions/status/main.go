package serverStatus

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/compute/v1"
)

var (
	projectID            string
	zone                 string
	instance             string
	userName             string
	password             string
	telnetHost           string
	telnetPort           string
	telnetPassword       string
	computeServiceClient *compute.Service
)

func init() {
	log.Println("Initializing server status function")

	// Get env vars
	projectID = os.Getenv("GCP_PROJECT_ID")
	zone = os.Getenv("GCP_ZONE")
	instance = os.Getenv("GCP_INSTANCE_NAME")
	userName = os.Getenv("USER_NAME")
	passwordSecretName := os.Getenv("PASSWORD_SEC_NAME")
	telnetHost = os.Getenv("TELNET_HOST")
	telnetPort = os.Getenv("TELNET_PORT")
	telnetPasswordSecName := os.Getenv("TELNET_PASSWORD_SEC_NAME")

	if projectID == "" ||
		zone == "" ||
		instance == "" ||
		userName == "" ||
		passwordSecretName == "" ||
		telnetHost == "" ||
		telnetPort == "" ||
		telnetPasswordSecName == "" {
		log.Fatalf(
			"Failed to get required environment variables:"+
				"projectID: '%s',"+
				"zone: '%s',"+
				"instance: '%s',"+
				"userName: '%s',"+
				"passwordSecretName: '%s',"+
				"telnetHost: '%s',"+
				"telnetPort: '%s',"+
				"telnetPasswordSecName: '%s'",
			projectID,
			zone,
			instance,
			userName,
			passwordSecretName,
			telnetHost,
			telnetPort,
			telnetPasswordSecName,
		)
	}

	ctx := context.Background()

	// Create GCP Secret Manager client
	secretManagerClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to setup Secret Manager client: %v", err)
	}
	defer secretManagerClient.Close()

	// Get password secret
	reqSecPassword := &secretmanagerpb.AccessSecretVersionRequest{Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, passwordSecretName)}
	respSecPassword, err := secretManagerClient.AccessSecretVersion(ctx, reqSecPassword)
	if err != nil {
		log.Fatalf("Failed to get secret 'projects/%s/secrets/%s/versions/latest' from Secret Manager: %v", projectID, passwordSecretName, err)
	}

	password = string(respSecPassword.Payload.Data)

	// Get telnet password secret
	reqSecTelnetPassword := &secretmanagerpb.AccessSecretVersionRequest{Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, telnetPasswordSecName)}
	respSecTelnetPassword, err := secretManagerClient.AccessSecretVersion(ctx, reqSecTelnetPassword)
	if err != nil {
		log.Fatalf("Failed to get secret 'projects/%s/secrets/%s/versions/latest' from Secret Manager: %v", projectID, passwordSecretName, err)
	}

	telnetPassword = string(respSecTelnetPassword.Payload.Data)

	// Create GCP Compute Engine client
	computeServiceClient, err = compute.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to setup Compute Engine client: %v", err)
	}

	log.Println("Done initializing server status function")
}

// Entrypoint is the starting point for this Cloud Function
func Entrypoint(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving new request")
	basicAuth(getServerStatus, userName, password).ServeHTTP(w, r)
	return
}

// Check if VM is running and how many players are in the game (implicitely checks if game server is running as well)
func getServerStatus(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting server status")

	isRunning := isVMRunning(projectID, zone, instance)
	if !isRunning {
		fmt.Fprintln(w, "Server virtual machine is NOT running")
		return
	}

	players, err := getPlayersInGame(telnetHost, telnetPort, telnetPassword)
	if err != nil {
		fmt.Fprintln(w, "Server virtual machine is running\nError connecting to game telnet server for checking players in game, please retry in a few seconds")
		return
	}

	fmt.Fprintf(w, "Server virtual machine is running\nPlayers in game: %s", players)
}

// Check if VM status is 'RUNNING'
func isVMRunning(projectID, zone, instance string) bool {
	log.Printf("Checking VM status")
	vm, err := computeServiceClient.Instances.Get(projectID, zone, instance).Do()
	if err != nil {
		log.Printf("Error checking VM status: %v", err)
		return false
	}

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
