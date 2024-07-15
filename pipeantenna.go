package main

import (
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	var l sync.Mutex

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		l.Lock()
		w.Header().Set("Content-Disposition",
			"attachment; filename=\""+time.Now().Format("20060102_150405")+"\"")
		w.Header().Set("Content-Type", "application/octet-stream")
		io.Copy(w, os.Stdin)
		time.AfterFunc(3*time.Second, func() {
			os.Exit(0)
		})
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		l.Lock()
		io.Copy(os.Stdout, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
		time.AfterFunc(3*time.Second, func() {
			os.Exit(0)
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		d, _ := os.ReadFile("web_pages/index.html")
		w.Write(d)
	})

	http.ListenAndServeTLS(":8080", os.Args[len(os.Args)-2], os.Args[len(os.Args)-1], nil)
}
