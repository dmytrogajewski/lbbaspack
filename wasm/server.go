package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := ":8000"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	// Create a custom file server that sets correct MIME types
	fileServer := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set correct MIME types for WebAssembly files
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		// Disable MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		fileServer.ServeHTTP(w, r)
	})

	fmt.Printf("WebAssembly server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	log.Fatal(http.ListenAndServe(port, nil))
}
