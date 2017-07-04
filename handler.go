package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	//key file
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error1."))
		return
	}
	defer file.Close()
	//file type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error2."))
		return
	}
	fileType := http.DetectContentType(buff)
	fmt.Println(fileType)

	if checkFileType(fileType) == false {
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

	path := GetPathByMd5(imageId)
	//path := "E:/1037u/1.gif"

	if FileExist(path) {
		w.Write([]byte(imageId))
		return
	}

	err = MkdirByMd5(imageId)
	if err != nil {
		log.Println(err)
	}

	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}

	fmt.Println(path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0775)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error1."))
		return
	}
	defer f.Close()
	bytesWritten, err := io.Copy(f, file)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error2."))
		return
	}
	fmt.Println(bytesWritten)
	w.Write([]byte(imageId))

}

func UrlHandler(w http.ResponseWriter, r *http.Request) {

	//key url
	r.ParseForm()
	url := r.FormValue("url")
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		w.Write([]byte("Error:Image Download Error."))
	}

	defer resp.Body.Close()

	fmt.Println(resp)
	if resp.StatusCode == 200 {

		urlPath := strings.Split(url, "/")
		var name string
		if len(urlPath) > 1 {
			name = urlPath[len(urlPath)-1]
		}
		random := strconv.Itoa(rand.Int())
		name = conf.Storage + random + name
		fmt.Println(name)
		out, err := os.Create(name)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Upload Error0."))
			return
		}
		defer out.Close()

		pix, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Upload Error0."))
			return
		}
		_, err = io.Copy(out, bytes.NewReader(pix))
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Upload Error0."))
			return
		}
		file, err := os.Open(name)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Upload Error1."))
			return
		}
		defer file.Close()
		//file type
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Upload Error2."))
			return
		}
		fileType := http.DetectContentType(buff)
		fmt.Println(fileType)

		if checkFileType(fileType) == false {
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

		path := GetPathByMd5(imageId)
		//path := "E:/1037u/1.gif"

		err = MkdirByMd5(imageId)
		if err != nil {
			log.Println(err)
		}

		if FileExist(path) {
			w.Write([]byte(imageId))
			return
		}

		if _, err = file.Seek(0, 0); err != nil {
			log.Println(err)
		}

		fmt.Println(path)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0775)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Save Error1."))
			return
		}
		defer f.Close()
		bytesWritten, err := io.Copy(f, file)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error:Save Error2."))
			return
		}
		fmt.Println(bytesWritten)
		w.Write([]byte(imageId))

	} else {
		w.Write([]byte("Error:Remote Image Download Error."))
	}

	return

}

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	//key base64
	r.ParseForm()
	base64String := r.FormValue("base64")
	fmt.Println(base64String)
	base64String = deleteBase64Head(base64String)
	base64DecodeString, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		w.Write([]byte("Error:Base64 Decode Error."))
		return
	}

	random := strconv.Itoa(rand.Int())
	fmt.Println(random)
	name := random + "base64.jpg"
	name = conf.Storage + name
	err = ioutil.WriteFile(name, base64DecodeString, 0755)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error0."))
		return
	}

	file, err := os.Open(name)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error1."))
		return
	}
	defer file.Close()

	//file type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error2."))
		return
	}
	fileType := http.DetectContentType(buff)
	fmt.Println(fileType)

	if checkFileType(fileType) == false {
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

	path := GetPathByMd5(imageId)
	//path := "E:/1037u/1.gif"

	err = MkdirByMd5(imageId)
	if err != nil {
		log.Println(err)
	}

	if FileExist(path) {
		w.Write([]byte(imageId))
		return
	}

	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}

	fmt.Println(path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0775)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error1."))
		return
	}
	defer f.Close()
	bytesWritten, err := io.Copy(f, file)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error2."))
		return
	}
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

	//image format equal jpeg
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
	w.Write([]byte("<html><body><center><h1>It Works!</h1></center><hr><center>goImg</center></body></html>"))
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

func checkFileType(fileType string) bool {
	fileTypes := []string{"image/jpeg", "image/gif", "image/png", "image/webp"}
	result := false
	for _, v := range fileTypes {
		if fileType == v {
			result = true
			break
		}
	}
	return result
}
func deleteBase64Head(base64String string) string {
	deleteStrings := []string{"data:image/jpeg;base64,", "data:image/gif;base64,", "data:image/png;base64,", "data:image/webp;base64,"}
	for _, v := range deleteStrings {
		if strings.Contains(base64String, v) {
			base64String = strings.Replace(base64String, v, "", -1)
			break
		}

	}
	return base64String
}
