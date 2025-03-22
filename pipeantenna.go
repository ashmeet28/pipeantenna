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
	if len(os.Args) != 6 {
		os.Exit(1)
	}

	indexFilePath := os.Args[1]
	indexFileData, err := os.ReadFile(indexFilePath)
	if err != nil {
		os.Exit(1)
	}

	authKeyBuf := make([]byte, 4)

	if _, err := rand.Read(authKeyBuf); err != nil {
		os.Exit(1)
	}

	authKey := hex.EncodeToString(authKeyBuf)

	authKeyFilePath := os.Args[4]
	if err := os.WriteFile(authKeyFilePath, append([]byte(authKey), 0x0a), 0600); err != nil {
		os.Exit(1)
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

	certFilePath := os.Args[2]
	keyFilePath := os.Args[3]
	httpPortToListenOn := os.Args[5]

	http.ListenAndServeTLS(":"+httpPortToListenOn, certFilePath, keyFilePath, nil)
}
