package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./uploads/"

func main() {
	// Create upload directory if not exists
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, os.ModePerm)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)

	fmt.Println("Starting server at :8080")
	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8090", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
		<body>
		<h1>Upload File</h1>
		<form enctype="multipart/form-data" action="/upload" method="post">
		  <input type="file" name="file" />
		  <input type="submit" value="upload" />
		</form>
		</body>
		</html>
	`))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the uploaded file
	r.ParseMultipartForm(10 << 20) // 10MB max size
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error parsing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the file on the server
	filePath := filepath.Join(uploadPath, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the file data
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Extract file name from URL
	fileName := r.URL.Path[len("/download/"):]
	filePath := filepath.Join(uploadPath, fileName)

	// Serve the file for download
	http.ServeFile(w, r, filePath)
}
