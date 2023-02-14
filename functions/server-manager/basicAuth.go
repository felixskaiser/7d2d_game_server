package serverManager

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net"
	"net/http"
)

// Basic auth middleware
func basicAuth(next http.HandlerFunc, expectedUsername, expectedPassword string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authenticating request")
		reqInfo, err := getReqInfo(r)
		if err != nil {
			log.Printf("Couldn't get requester info from request: %s", err)
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}

		username, password, ok := r.BasicAuth()

		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
			expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				log.Printf("Authenticated request from %s", reqInfo)
				next.ServeHTTP(w, r)
				return
			}
		}

		log.Printf("Unauthenticated request from %s", reqInfo)
		w.Header().Set("WWW-Authenticate", `Basic realm="7d2dgameserver", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// Get some basic info about requester for logging
func getReqInfo(r *http.Request) (string, error) {
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	ipParsed := net.ParseIP(ip)
	forwarded := r.Header.Get("X-Forwarded-For")

	return fmt.Sprintf("remote address: '%s:%s', forwarded: %s", ipParsed, port, forwarded), nil
}
