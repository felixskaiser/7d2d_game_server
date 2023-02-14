package serverManager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"google.golang.org/api/compute/v1"
)

func startServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving '/start'")

	vm, err := computeServiceClient.Instances.Get(projectID, zone, instance).Do()
	if err != nil {
		log.Printf("SERVER_START: Error checking VM status: %v", err)
		startServerRespond(w, "error getting info about server")
		return
	}

	log.Printf("SERVER_START: VM status: %v", vm.Status)

	if vm.Status != "RUNNING" {
		err := startComputeInstance(ctx, projectID, zone, instance)
		if err != nil {
			log.Printf("SERVER_START: Error starting VM: %v", err)
			startServerRespond(w, "error starting server")
			return
		}
	} else {
		log.Println("SERVER_START: VM is already running")
		startServerRespond(w, "server is already running")
		return
	}

	startServerRespond(w, "started server")
	return
}

func startServerRespond(w http.ResponseWriter, response string) {
	layoutTmplPath := filepath.Join(templateDir, "base.tmpl")
	startTmplPath := filepath.Join(templateDir, "start.tmpl")

	tmpl, err := template.ParseFiles(layoutTmplPath, startTmplPath)
	if err != nil {
		log.Printf("SERVER_START: Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", response)
	if err != nil {
		log.Printf("SERVER_START: Error serving '/start': %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	return
}

func startComputeInstance(ctx context.Context, projectID, zone, instance string) error {
	op, err := computeServiceClient.Instances.Start(projectID, zone, instance).Do()
	if err != nil {
		return fmt.Errorf("Error starting VM: %v", err)
	}

	err = waitForOperation(ctx, projectID, zone, op)
	if err != nil {
		return fmt.Errorf("Error waiting for VM to start: %v", err)
	}

	return nil
}

func waitForOperation(ctx context.Context, project, zone string, op *compute.Operation) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for operation to complete")
		case <-ticker.C:
			result, err := computeServiceClient.ZoneOperations.Get(project, zone, op.Name).Do()
			if err != nil {
				return fmt.Errorf("ZoneOperations.Get: %s", err)
			}

			if result.Status == "DONE" {
				if result.Error != nil {
					var errors []string
					for _, e := range result.Error.Errors {
						errors = append(errors, e.Message)
					}
					return fmt.Errorf("operation %q failed with error(s): %s", op.Name, strings.Join(errors, ", "))
				}
				return nil
			}
		}
	}
}
