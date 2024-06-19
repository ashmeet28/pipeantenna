package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+os.Args[1]+"\"")
		w.Header().Set("Content-Type", "application/octet-stream")

		io.Copy(w, os.Stdin)
		go func() {
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}()
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(os.Stdout, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
		go func() {
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>File Upload</title>
				<script>
					function uploadFile(fileInput) {
						const file = fileInput.files[0];
						if (!file) {
							alert('Please select a file');
							return;
						}

						alert("Starting file upload")

						const xhr = new XMLHttpRequest();
						xhr.open('POST', '/upload', true);
						xhr.setRequestHeader('Content-Type', file.type);
						xhr.onload = function () {
							if (xhr.status === 200) {
								alert(xhr.responseText);
							}
						};
						xhr.send(file);
					}
				</script>
			</head>
			<body>
				<h1>File Upload</h1>
				<input type="file" id="fileInput" />
				<button onclick="uploadFile(document.getElementById('fileInput'))">Upload File</button>
				<h1>File Download</h1>
				<form action="/download" method="get">
					<button type="submit">Download File</button>
				</form>
			</body>
			</html>
		`
		io.WriteString(w, html)
	})

	http.ListenAndServe(":8080", nil)
}
