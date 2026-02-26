package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"php-dep-extractor/internal/server"
)

//go:embed web/*
var webFS embed.FS

const Version = "0.1.0"

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	url := fmt.Sprintf("http://127.0.0.1:%d", addr.Port)
	fmt.Printf("PHP Dependency Extractor v%s running at %s\n", Version, url)

	srv := server.New(webFS)

	go openBrowser(url)

	if err := http.Serve(listener, srv); err != nil {
		log.Fatal("Server error:", err)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
