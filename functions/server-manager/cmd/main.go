// The cmd command starts an HTTP server for developing the GCP Cloud Function locally
package main

import (
	"fmt"
	"log"
	"net/http"

	serverManager "github.com/felixskaiser/7d2d_game_server/server-manager"
)

func main() {
	http.HandleFunc("/", serverManager.Entrypoint)
	fmt.Println("Listening on localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
