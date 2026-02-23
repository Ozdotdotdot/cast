package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

// Start creates and starts an HTTP server to serve the audio file
// Returns the server, the URL to access the file, and any error
func Start(audioPath string) (*Server, string, error) {
	// Get local IP
	localIP, err := getLocalIP()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get local IP: %w", err)
	}

	// Find a free port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, "", fmt.Errorf("failed to find free port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	filename := filepath.Base(audioPath)

	// Create file server
	mux := http.NewServeMux()
	mux.HandleFunc("/"+filename, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/mpeg")
		http.ServeFile(w, r, audioPath)
	})

	httpServer := &http.Server{
		Handler: mux,
	}

	// Start server in background
	go func() {
		httpServer.Serve(listener)
	}()

	url := fmt.Sprintf("http://%s:%d/%s", localIP, port, filename)

	return &Server{
		httpServer: httpServer,
		listener:   listener,
	}, url, nil
}

// Shutdown gracefully stops the HTTP server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

// getLocalIP returns the local IP address used to reach the internet
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
