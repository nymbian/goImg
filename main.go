package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("start")
	//test()
	LoadConf()
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/upload", UploadHandler).Methods("POST")
	r.HandleFunc("/url", UrlHandler).Methods("POST")
	r.HandleFunc("/base64", Base64Handler).Methods("POST")
	r.HandleFunc("/{imageId}_{width:[0-9]+}x{height:[0-9]+}", DownloadResizHandler).Methods("GET")
	r.HandleFunc("/{imageId}", DownloadHandler).Methods("GET")
	err := http.ListenAndServe(conf.ListenAddr, r)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
