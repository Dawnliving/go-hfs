package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type JsonResult struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
}


func uploadFile(w http.ResponseWriter, r *http.Request) {
	msg, _ := json.Marshal(JsonResult{Code: 404, Msg: "upload error-occur"})
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.WriteHeader(404)
		w.Write(msg)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	tempFile, err := os.CreateTemp("/fordev/upload/image", "IMG_*.png")//linux env
	//tempFile, err := os.CreateTemp("D:/golang_projects/go_http_file_system/test", "IMG_*.png")//windows test env
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		w.Write(msg)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// return that we have successfully uploaded our file!
	fmt.Printf("successfully uploaded")

	res := filepath.Base(tempFile.Name())
	if res != "" {
		msg, _ = json.Marshal(JsonResult{Code: 200, Msg: res})
		w.WriteHeader(200)
		w.Write(msg)
	}


}

func setupRoutes() {
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("/fordev/upload/image"))))//linux env
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("D:/golang_projects/go_http_file_system/test"))))//windows env
	http.HandleFunc("/upload", uploadFile)
	err := http.ListenAndServe(":16535", nil)
	if err != nil{
		log.Fatal("ListenAndServe", err)
	}
}


func main()  {
	setupRoutes()
}
