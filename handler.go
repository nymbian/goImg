package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	//key uploadfile
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error.1"))
		return
	}
	defer file.Close()
	//file type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error.2"))
		return
	}
	fileType := http.DetectContentType(buff)
	fmt.Println(fileType)

	fileTypes := []string{"image/jpeg", "image/gif", "image/png", "image/webp"}
	result := false
	for _, v := range fileTypes {
		if fileType == v {
			result = true
			break
		}
	}

	if result == false {
		w.Write([]byte("Error:Not Img."))
		return
	}

	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}

	md5h := md5.New()
	io.Copy(md5h, file)

	imageId := hex.EncodeToString(md5h.Sum(nil))
	fmt.Println(imageId)
	//imageId = "86427d1debefe65f0da3a7affdc204f2"

	err = MkdirByMd5(imageId)
	if err != nil {
		log.Println(err)
	}

	path := GetPathByMd5(imageId)
	//path := "E:/1037u/1.gif"

	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}

	fmt.Println(path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error."))
		return
	}
	defer f.Close()
	bytesWritten, err := io.Copy(f, file)
	checkErr(err)
	fmt.Println(bytesWritten)
	w.Write([]byte(imageId))

}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageId := vars["imageId"]
	if len([]rune(imageId)) != 32 {
		w.Write([]byte("Error:ImageID incorrect."))
		return
	}
	imgPath := GetPathByMd5(imageId)
	if !FileExist(imgPath) {
		w.Write([]byte("Error:Image Not Found."))
		return
	}
	http.ServeFile(w, r, imgPath)
}

func DownloadResizHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageId := vars["imageId"]
	width := vars["width"]
	height := vars["height"]

	widthInt, _ := strconv.Atoi(width)
	heightInt, _ := strconv.Atoi(height)

	if len([]rune(imageId)) != 32 {
		w.Write([]byte("Error:ImageID incorrect."))
		return
	}
	imgPath := GetPathByMd5(imageId)

	file, err := os.Open(imgPath)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Image Open Error."))
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Image Not Found."))
		return
	}
	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}
	fileType := http.DetectContentType(buff)
	fmt.Println(fileType)

	if fileType == "image/jpeg" {
		resizeImgPath := imgPath + "_" + width + "x" + height

		// resize img not exist
		if !FileExist(resizeImgPath) {

			// decode jpeg into image.Image
			img, err := jpeg.Decode(file)
			if err != nil {
				w.Write([]byte("Error:Image Decode Error."))
				return
			}
			file.Close()

			// resize to width 1000 using Lanczos resampling
			// and preserve aspect ratio

			m := resize.Resize(uint(widthInt), uint(heightInt), img, resize.Lanczos3)

			out, err := os.Create(resizeImgPath)
			if err != nil {
				w.Write([]byte("Error:Resize Image Create Error."))
				return
			}
			defer out.Close()

			// write new image to file
			jpeg.Encode(out, m, nil)
		}

		http.ServeFile(w, r, resizeImgPath)

	} else {
		http.ServeFile(w, r, imgPath)
	}

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><center><h1>It Works!</h1></center><hr><center>Quick Image Server</center></body></html>"))
}

func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}
