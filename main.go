package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//test()
	LoadConf()
	log.Println("listen:" + conf.ListenAddr)

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")                                                  //home
	r.HandleFunc("/upload", UploadHandler).Methods("POST")                                         //upload file
	r.HandleFunc("/url", UrlHandler).Methods("POST")                                               //url url
	r.HandleFunc("/base64", Base64Handler).Methods("POST")                                         //base64 base64
	r.HandleFunc("/{imageId}_{width:[0-9]+}x{height:[0-9]+}", DownloadResizHandler).Methods("GET") //resize
	r.HandleFunc("/{imageId}_{action:sync}", DownloadHandler).Methods("GET")                       //sync
	r.HandleFunc("/{imageId}", DownloadHandler).Methods("GET")                                     //normal

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
