package serverManager

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

const (
	templateSubDir string = "templates/"
)

var (
	templateDir          string
	projectID            string
	zone                 string
	instance             string
	userName             string
	password             string
	telnetHost           string
	telnetPort           string
	telnetPassword       string
	computeServiceClient *compute.Service

	mux = http.NewServeMux()
)

type Empty struct{}

func init() {
	log.Println("Initializing server status function")

	// Get env vars
	iscloudFunctionStr := os.Getenv("IS_CLOUD_FUNCTION")
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

	if iscloudFunctionStr == "true" {
		// DEBUG
		_, err := os.Stat("/workspace")
		if err != nil {
			log.Printf("Not a directory: /workspace': %v", err)
		} else {
			entries, _ := os.ReadDir("/workspace")
			log.Printf("DEBUG: '/workspace' => %v", entries)
			for _, e := range entries {
				log.Printf("DEBUG: %v", e.Name())
			}
		}

		_, err = os.Stat("/workspace/serverless_function_source_code")
		if err != nil {
			log.Printf("Not a directory: ./serverless_function_source_code': %v", err)
		} else {
			entries, _ := os.ReadDir("./serverless_function_source_code")
			log.Printf("DEBUG: './serverless_function_source_code' => %v", entries)
			for _, e := range entries {
				log.Printf("DEBUG: %v", e.Name())
			}
		}

		functionSourceCodeDir := "/workspace/serverless_function_source_code/"
		_, err = os.Stat(functionSourceCodeDir)
		if err != nil {
			log.Fatalf("Not a directory: '%s': %v", functionSourceCodeDir, err)
		}

		templateDir = functionSourceCodeDir + templateSubDir
		_, err = os.Stat(templateDir)
		if err != nil {
			log.Fatalf("Not a directory: '%s': %v", functionSourceCodeDir, err)
		}

	} else {
		templateDir = templateSubDir
	}
	log.Printf("Looking for HTML template files at %s", templateDir)

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

	// Register routes
	mux.HandleFunc("/status", basicAuth(getServerStatus, userName, password))
	mux.HandleFunc("/", basicAuth(defaultHandler, userName, password))

	log.Println("Done initializing server status function")
}

// Entrypoint is the starting point for this Cloud Function
func Entrypoint(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
	return
}
