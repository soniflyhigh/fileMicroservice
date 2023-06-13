package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	templateFile = template.Must(template.ParseFiles("template/index.html"))
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleUpload(w, r)
		return
	}
	templateFile.ExecuteTemplate(w, "index.html", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, fileHeader, err := r.FormFile("image")
	fmt.Println("FileName", fileHeader.Filename)
	fmt.Println("FileSize", fileHeader.Size)
	fmt.Println("FileHeader", fileHeader.Header)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("", "upload-*.tmp")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	buffer := make([]byte, 1<<20)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println("this is value of n", n)
		tempFile.Write(buffer[:n])
	}

	fmt.Println("this is tempriory file", tempFile.Name())

	filename := path.Base(fileHeader.Filename)
	fmt.Println(tempFile.Name())
	dest, err := os.Create(filename)
	fmt.Println("destination",&dest)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?sucess=true", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", home)
	log.Println("Starting server port on 5007:")
	if err := http.ListenAndServe(":5007", nil); err != nil {
		log.Fatal(err)
	}
}
