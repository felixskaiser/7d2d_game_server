package serverStatus

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	// Some typical errors one runs into when connecting to game's telnet server
	telnetErrors = []string{
		"ERR Exception in thread",
		"EXC The operation is not allowed on non-connected sockets",
		"EXC Unable to write data to the transport connection",
	}
)

// Send command to game server via game's telnet server
func sendTelnetCmd(cmd, host, port, password string) (string, error) {
	// Connect to game's telnet server via TCP
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return "", fmt.Errorf("TELNET_SERVER_ERROR: Error connecting to telnet server: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	defer conn.Close()

	log.Printf("TELNET_SERVER: Connection established to telnet server (%s)", conn.RemoteAddr().String())

	err = login(password, conn)
	if err != nil {
		return "", fmt.Errorf("TELNET_SERVER_ERROR: Error logging in to telnet server: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	resp, err := sendCmd(cmd, conn)
	if err != nil {
		return "", fmt.Errorf("TELNET_SERVER: Error sending command to telnet server: %v", err)
	}

	return resp, nil
}

// Login to game's telnet server
func login(password string, conn net.Conn) error {
	writerLogin := bufio.NewWriter(conn)
	readerLogin := bufio.NewReader(conn)

	// Expect prompt to login
	respLoginInit, err := readerLogin.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Error initializing reader for login to telnet server: %v", err)
	}
	log.Printf("TELNET_SERVER_LOGIN: Initializing reader for login to telnet server, got response: '%s'", respLoginInit)

	// Send password to login
	_, err = writerLogin.WriteString(fmt.Sprintf("%s\n", password))
	if err != nil {
		return fmt.Errorf("Error writing password to login for telnet server: %v", err)
	}

	err = writerLogin.Flush()
	if err != nil {
		return fmt.Errorf("Error logging in to telnet server: %v", err)
	}

	// Listen for reply
	respLogin, err := readerLogin.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Error logging in to telnet server: %v", err)
	}
	log.Printf("TELNET_SERVER_LOGIN: Logging in to telnet server, got response: %s", respLogin)

	writerLogin.Reset(conn)
	readerLogin.Reset(conn)
	return nil
}

// Send actual command to game's telnet server. Requires previous login
func sendCmd(cmd string, conn net.Conn) (string, error) {
	writerCmd := bufio.NewWriter(conn)
	readerCmd := bufio.NewReader(conn)

	log.Printf("TELNET_SERVER_CMD: Sending command '%s' to telnet server", cmd)

	// Read until no more content
	err := readAll(conn, readerCmd, 5)
	if err != nil {
		return "", err
	}

	// Send command to telnet server
	_, err = writerCmd.WriteString(fmt.Sprintf("%s\n", cmd))
	if err != nil {
		return "", fmt.Errorf("Error writing command for telnet server: %v", err)
	}

	err = writerCmd.Flush()
	if err != nil {
		return "", fmt.Errorf("Error sending command '%s' to telnet server: %v", cmd, err)
	}

	// Listen for reply
	respCmd, err := readerCmd.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Error sending command '%s' to telnet server: %v", cmd, err)
	}
	log.Printf("TELNET_SERVER_CMD: Sent command '%s' to telnet server, got response: %s", cmd, respCmd)

	err = filterRespForKnownErrors(respCmd, telnetErrors)
	if err != nil {
		return "", err
	}

	// Read until no more content
	err = readAll(conn, readerCmd, 5)
	if err != nil {
		return "", err
	}

	return respCmd, nil
}

// Attempt to read all content from game's telnet server to clear the reader but still be able to log server's responses
func readAll(conn net.Conn, reader *bufio.Reader, maxTries int) error {
	tries := 0
	for {
		resp, err := reader.ReadString('\n')
		if err != nil {
			// Treat timeout as telnet server having nothing more to send
			if errors.Is(err, os.ErrDeadlineExceeded) {
				conn.SetReadDeadline(time.Now().Add(2 * time.Second))
				return nil
			}

			return fmt.Errorf("Error reading all from telnet server: %v", err)
		}

		if resp == "" {
			return nil
		}

		log.Printf("Reading all from telnet server, got response: %s", resp)

		if tries >= maxTries {
			return fmt.Errorf("Max tries for reading all from telnet server")
		}

		tries += 1
	}
}

// Attempt to filter actual errors in game's telnet server's response before responding to requester
func filterRespForKnownErrors(resp string, errors []string) error {
	for _, err := range errors {
		if strings.Contains(resp, err) {
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}
