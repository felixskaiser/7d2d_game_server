// The cmd command starts an HTTP server for developing the GCP Cloud Function locally
package main

import (
	"fmt"
	"log"
	"net/http"

	serverStatus "github.com/felixskaiser/7d2d_game_server/status"
)

func main() {
	http.HandleFunc("/", serverStatus.Entrypoint)
	fmt.Println("Listening on localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
