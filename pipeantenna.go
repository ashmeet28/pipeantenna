package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		return
	}

	indexFileData, err1 := os.ReadFile("web_pages/index.html")
	if err1 != nil {
		return
	}

	authKeyBuf := make([]byte, 4)

	_, err2 := rand.Read(authKeyBuf)
	if err2 != nil {
		return
	}

	authKey := hex.EncodeToString(authKeyBuf)

	err3 := os.WriteFile(os.Args[3], append([]byte(authKey), 0x0a), 0600)
	if err3 != nil {
		return
	}

	var l sync.Mutex

	http.HandleFunc("GET /download", func(w http.ResponseWriter, r *http.Request) {
		l.Lock()

		if r.Method == http.MethodGet && r.URL.Path == "/download" &&
			len(r.URL.Query()["auth_key"]) == 1 && r.URL.Query()["auth_key"][0] == authKey {

			w.Header().Set("Content-Disposition",
				"attachment; filename=\""+time.Now().Format("20060102_150405")+"\"")
			w.Header().Set("Content-Type", "application/octet-stream")
			io.Copy(w, os.Stdin)

		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		time.AfterFunc(3*time.Second, func() {
			os.Exit(0)
		})
	})

	http.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		l.Lock()

		if r.URL.Path == "/upload" {
			_, err := io.Copy(os.Stdout, r.Body)
			if err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

		time.AfterFunc(3*time.Second, func() {
			os.Exit(0)
		})
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexFileData)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.ListenAndServeTLS(":8080", os.Args[1], os.Args[2], nil)
}
